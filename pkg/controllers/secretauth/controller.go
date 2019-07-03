package secretauth

import (
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	nebulav1 "github.com/puppetlabs/nebula-tasks/pkg/apis/nebula.puppet.com/v1"
	"github.com/puppetlabs/nebula-tasks/pkg/config"
	"github.com/puppetlabs/nebula-tasks/pkg/data/secrets/vault"
	clientset "github.com/puppetlabs/nebula-tasks/pkg/generated/clientset/versioned"
	informers "github.com/puppetlabs/nebula-tasks/pkg/generated/informers/externalversions"
	sainformers "github.com/puppetlabs/nebula-tasks/pkg/generated/informers/externalversions/nebula.puppet.com/v1"
)

const (
	// default image:tag to use for nebula-metadata-api
	metadataServiceImage = "pcr-internal.puppet.net/nebula/nebula-metadata-api:latest"
	// default name for the workflow metadata api pod and service
	metadataServiceName = "workflow-metadata-api"
	// default maximum retry attempts to create resources spawned from SecretAuth creations
	maxRetries = 10
)

// Controller watches for nebulav1.SecretAuth resource changes.
// If a SecretAuth resource is created, the controller will create a service acccount + rbac
// for the namespace, then inform vault that that service account is allowed to access
// readonly secrets under a preconfigured path related to a nebula workflow. It will then
// spin up a pod running an instance of nebula-metadata-api that knows how to
// ask kubernetes for the service account token, that it will use to proxy secrets
// between the task pods and the vault server.
type Controller struct {
	kubeclient       kubernetes.Interface
	clientset        clientset.Interface
	workqueue        workqueue.RateLimitingInterface
	saInformer       sainformers.SecretAuthInformer
	saInformerSynced cache.InformerSynced
	informerFactory  informers.SharedInformerFactory
	vaultClient      *vault.VaultAuth
}

// Run starts all required informers and spawns two worker goroutines
// that will pull resource objects off the workqueue. This method blocks
// until stopCh is closed or an earlier bootstrap call results in an error.
func (c *Controller) Run(stopCh chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	c.informerFactory.Start(stopCh)

	if ok := cache.WaitForCacheSync(stopCh, c.saInformerSynced); !ok {
		return fmt.Errorf("failed to wait for informer cache to sync")
	}

	go wait.Until(c.worker, time.Second, stopCh)
	go wait.Until(c.worker, time.Second, stopCh)

	<-stopCh

	return nil
}

func (c *Controller) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	defer c.workqueue.Done(key)

	err := c.processSingleItem(key.(string))
	c.handleErr(err, key)

	return true
}

// handleErr takes an error and the k8s object key that caused it
// and tries to requeue the object for processing. If we have requeued
// that key the maximum number of times, then we drop the object from the
// queue and tell kubernetes runtime to handle the error.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.workqueue.Forget(key)
		return
	}

	// requeue if we can still retry to process the resource
	if c.workqueue.NumRequeues(key) < maxRetries {
		log.Printf("error syncing secretauth %v: %v", key, err)
		c.workqueue.AddRateLimited(key)
		return
	}

	// otherwise report the error to kubernetes and drop the key from the queue
	log.Printf("dropping secretauth from queue %v: %v", key, err)
	utilruntime.HandleError(err)
	c.workqueue.Forget(key)
}

// processSingleItem is responsible for creating all the resouces required for
// secret handling and authentication.
// TODO break this logic out into smaller chunks... especially the calls to the vault api
func (c *Controller) processSingleItem(key string) error {
	log.Println("syncing SecretAuth", key)
	defer log.Println("done syncing SecretAuth", key)

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	sa, err := c.clientset.NebulaV1().SecretAuths(namespace).Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if sa.Status.ServiceAccount != "" {
		log.Printf("resources for %s have already been created", key)
		return nil
	}

	var (
		saccount  *corev1.ServiceAccount
		pod       *corev1.Pod
		service   *corev1.Service
		configMap *corev1.ConfigMap
	)

	log.Println("creating service account for", sa.Spec.WorkflowID)
	saccount, err = c.kubeclient.CoreV1().ServiceAccounts(namespace).Create(serviceAccount(sa))
	if errors.IsAlreadyExists(err) {
		saccount, err = c.kubeclient.CoreV1().ServiceAccounts(namespace).Get(getName(sa), metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	log.Println("creating metadata service pod for", sa.Spec.WorkflowID)
	pod, err = c.kubeclient.CoreV1().Pods(namespace).Create(metadataServicePod(
		saccount, sa, c.vaultClient.Address(), c.vaultClient.EngineMount()))
	if errors.IsAlreadyExists(err) {
		pod, err = c.kubeclient.CoreV1().Pods(namespace).Get(metadataServiceName, metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	log.Println("creating pod service for", sa.Spec.WorkflowID)
	service, err = c.kubeclient.CoreV1().Services(namespace).Create(metadataServiceService(sa))
	if errors.IsAlreadyExists(err) {
		service, err = c.kubeclient.CoreV1().Services(namespace).Get(metadataServiceName, metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	log.Println("creating config map for", sa.Spec.WorkflowID)
	configMap, err = c.kubeclient.CoreV1().ConfigMaps(namespace).Create(workflowConfigMap(sa))
	if errors.IsAlreadyExists(err) {
		configMap, err = c.kubeclient.CoreV1().ConfigMaps(namespace).Get(getName(sa), metav1.GetOptions{})
	}
	if err != nil {
		return err
	}

	log.Println("writing vault readonly access policy for ", sa.Spec.WorkflowID)
	// now we let vault know about the service account
	if err := c.vaultClient.WritePolicy(namespace, sa.Spec.WorkflowID); err != nil {
		return err
	}

	log.Println("enabling vault access for workflow service account for ", sa.Spec.WorkflowID)
	if err := c.vaultClient.WriteRole(namespace, saccount.GetName(), namespace); err != nil {
		return err
	}

	saCopy := sa.DeepCopy()
	saCopy.Status.MetadataServicePod, err = cache.MetaNamespaceKeyFunc(pod)
	if err != nil {
		return err
	}

	saCopy.Status.MetadataServiceService, err = cache.MetaNamespaceKeyFunc(service)
	if err != nil {
		return err
	}

	saCopy.Status.ServiceAccount, err = cache.MetaNamespaceKeyFunc(saccount)
	if err != nil {
		return err
	}

	saCopy.Status.ConfigMap, err = cache.MetaNamespaceKeyFunc(configMap)
	if err != nil {
		return err
	}

	saCopy.Status.VaultPolicy = namespace
	saCopy.Status.VaultAuthRole = namespace

	log.Println("updating secretauth resource status for ", sa.Spec.WorkflowID)
	saCopy, err = c.clientset.NebulaV1().SecretAuths(namespace).Update(saCopy)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) enqueueSecretAuth(obj interface{}) {
	sa := obj.(*nebulav1.SecretAuth)

	key, err := cache.MetaNamespaceKeyFunc(sa)
	if err != nil {
		utilruntime.HandleError(err)

		return
	}

	c.workqueue.Add(key)
}

func NewController(cfg *config.SecretAuthControllerConfig, vaultClient *vault.VaultAuth) (*Controller, error) {
	kcfg, err := clientcmd.BuildConfigFromFlags(cfg.KubeMasterURL, cfg.Kubeconfig)
	if err != nil {
		return nil, err
	}

	kc, err := kubernetes.NewForConfig(kcfg)
	if err != nil {
		return nil, err
	}

	nebclient, err := clientset.NewForConfig(kcfg)
	if err != nil {
		return nil, err
	}

	nebInformerFactory := informers.NewSharedInformerFactory(nebclient, time.Second*30)
	saInformer := nebInformerFactory.Nebula().V1().SecretAuths()

	c := &Controller{
		kubeclient:       kc,
		clientset:        nebclient,
		workqueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "SecretAuths"),
		saInformer:       saInformer,
		saInformerSynced: saInformer.Informer().HasSynced,
		informerFactory:  nebInformerFactory,
		vaultClient:      vaultClient,
	}

	saInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.enqueueSecretAuth,
	})

	return c, nil
}

func serviceAccount(sa *nebulav1.SecretAuth) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(sa),
			Namespace: sa.GetNamespace(),
		},
		ImagePullSecrets: []corev1.LocalObjectReference{
			{
				Name: "image-pull-secret",
			},
		},
	}
}

func metadataServicePod(saccount *corev1.ServiceAccount, sa *nebulav1.SecretAuth, vaultAddr, vaultEngineMount string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metadataServiceName,
			Namespace: sa.GetNamespace(),
			Labels: map[string]string{
				"app": metadataServiceName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            metadataServiceName,
					Image:           metadataServiceImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command: []string{
						"/usr/bin/nebula-metadata-api",
						"-bind-addr",
						":7000",
						"-vault-addr",
						defaultVaultAddr,
						"-vault-role",
						sa.GetNamespace(),
						"-workflow-name",
						sa.Spec.WorkflowID,
						"-vault-engine-mount",
						vaultEngineMount,
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 7000,
						},
					},
				},
			},
			ServiceAccountName: saccount.GetName(),
			RestartPolicy:      corev1.RestartPolicyOnFailure,
		},
	}
}

func metadataServiceService(sa *nebulav1.SecretAuth) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metadataServiceName,
			Namespace: sa.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(7000),
				},
			},
			Selector: map[string]string{
				"app": metadataServiceName,
			},
		},
	}
}

func workflowConfigMap(sa *nebulav1.SecretAuth) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(sa),
			Namespace: sa.GetNamespace(),
		},
		Data: map[string]string{
			"metadata-api-url": fmt.Sprintf("http://%s.%s", metadataServiceName, sa.GetNamespace()),
		},
	}
}

func getName(sa *nebulav1.SecretAuth) string {
	return fmt.Sprintf("%s-%d", sa.Spec.WorkflowID, sa.Spec.RunNum)
}
