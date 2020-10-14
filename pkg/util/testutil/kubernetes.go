package testutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/puppetlabs/relay-core/pkg/operator/dependency"
	"github.com/puppetlabs/relay-core/pkg/util/retry"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	tektonfake "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/fake"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
	cachingv1alpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	TestScheme = runtime.NewScheme()
)

func init() {
	schemeBuilder := runtime.NewSchemeBuilder(
		dependency.AddToScheme,
		apiextensionsv1.AddToScheme,
		apiextensionsv1beta1.AddToScheme,
		cachingv1alpha1.AddToScheme,
	)

	if err := schemeBuilder.AddToScheme(TestScheme); err != nil {
		panic(err)
	}
}

func NewMockKubernetesClient(initial ...runtime.Object) kubernetes.Interface {
	for _, obj := range initial {
		setObjectUIDOnObject(obj)
	}

	kc := fake.NewSimpleClientset(initial...)
	kc.PrependReactor("create", "*", setObjectUID)
	kc.PrependReactor("list", "pods", filterListPods(kc.Tracker()))
	return kc
}

func NewMockTektonKubernetesClient(initial ...runtime.Object) tekton.Interface {
	for _, obj := range initial {
		setObjectUIDOnObject(obj)
	}

	tkc := tektonfake.NewSimpleClientset(initial...)
	tkc.PrependReactor("create", "*", setObjectUID)
	return tkc
}

func setObjectUID(action testing.Action) (bool, runtime.Object, error) {
	switch action := action.(type) {
	case testing.CreateActionImpl:
		obj := action.GetObject()
		setObjectUIDOnObject(obj)
		return false, obj, nil
	default:
		return false, nil, nil
	}
}

func setObjectUIDOnObject(obj runtime.Object) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return
	}

	accessor.SetUID(types.UID(uuid.New().String()))
}

func filterListPods(tracker testing.ObjectTracker) testing.ReactionFunc {
	delegate := testing.ObjectReaction(tracker)

	return func(action testing.Action) (bool, runtime.Object, error) {
		la := action.(testing.ListAction)

		found, obj, err := delegate(action)
		if err != nil || !found {
			return found, obj, err
		}

		pods := obj.(*corev1.PodList)

		keep := 0
		for _, pod := range pods.Items {
			if !la.GetListRestrictions().Fields.Matches(fields.Set{
				"status.podIP": pod.Status.PodIP,
				"status.phase": string(pod.Status.Phase),
			}) {
				continue
			}

			pods.Items[keep] = pod
			keep++
		}

		pods.Items = pods.Items[:keep]
		return true, pods, nil
	}
}

type ParseKubernetesManifestPatcherFunc func(obj runtime.Object, gvk *schema.GroupVersionKind)

func ParseKubernetesManifest(r io.ReadCloser, patchers ...ParseKubernetesManifestPatcherFunc) ([]runtime.Object, error) {
	patchers = append(patchers, KubernetesFixupPatcher)

	decoder := yaml.NewDocumentDecoder(r)
	defer decoder.Close()

	// Copy buffer; we can't use io.Copy because of the weird semantics of the
	// document decoder in how it returns ErrShortBuffer.
	buf := make([]byte, 32*1024)

	// This lets us convert input documents.
	deserializer := serializer.NewCodecFactory(TestScheme).UniversalDeserializer()

	// The objects to create.
	var objs []runtime.Object

	var stop bool
	for !stop {
		var doc bytes.Buffer

		for {
			nr, err := decoder.Read(buf)
			if nr > 0 {
				if nw, err := doc.Write(buf[:nr]); err != nil {
					return nil, err
				} else if nw != nr {
					return nil, io.ErrShortWrite
				}
			}

			if err == io.ErrShortWrite {
				// More document to read, keep going.
			} else if err == io.EOF {
				// End of the entire stream.
				stop = true
				break
			} else if err != nil {
				return nil, err
			} else {
				// End of this loop, but we have another document ahead.
				break
			}
		}

		b := doc.Bytes()
		if len(bytes.TrimSpace(b)) == 0 {
			// Empty document.
			continue
		}

		obj, gvk, err := deserializer.Decode(b, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, patcher := range patchers {
			patcher(obj, gvk)
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

func KubernetesFixupPatcher(obj runtime.Object, gvk *schema.GroupVersionKind) {
	switch t := obj.(type) {
	case *appsv1.Deployment:
		// SSA has marked "protocol" is required but basically everyone expects
		// it to default to TCP.
		for i, container := range t.Spec.Template.Spec.Containers {
			for j, port := range container.Ports {
				if len(port.Protocol) > 0 {
					continue
				}

				t.Spec.Template.Spec.Containers[i].Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	case *corev1.Service:
		// Same for services.
		for i, port := range t.Spec.Ports {
			if len(port.Protocol) > 0 {
				continue
			}

			t.Spec.Ports[i].Protocol = corev1.ProtocolTCP
		}
	}
}

func KubernetesDefaultNamespacePatcher(mapper meta.RESTMapper, namespace string) ParseKubernetesManifestPatcherFunc {
	return func(obj runtime.Object, gvk *schema.GroupVersionKind) {
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return
		}

		// Namespace already set?
		if accessor.GetNamespace() != "" {
			return
		}

		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return
		}

		// Does this resource even take a namespace?
		if mapping.Scope.Name() != meta.RESTScopeNameNamespace {
			return
		}

		accessor.SetNamespace(namespace)
	}
}

func ParseApplyKubernetesManifest(ctx context.Context, cl client.Client, r io.ReadCloser, patchers ...ParseKubernetesManifestPatcherFunc) ([]runtime.Object, error) {
	objs, err := ParseKubernetesManifest(r, patchers...)
	if err != nil {
		return nil, err
	}

	for i, obj := range objs {
		md, err := meta.Accessor(obj)
		if err != nil {
			log.Printf("... applying %s (error retrieving object name: %+v)", obj.GetObjectKind().GroupVersionKind(), err)
		} else {
			name := md.GetName()
			if md.GetNamespace() != "" {
				name = md.GetNamespace() + "/" + name
			}
			log.Printf("... applying %s %s", obj.GetObjectKind().GroupVersionKind(), name)
		}

		if err := cl.Patch(ctx, obj, client.Apply, client.ForceOwnership, client.FieldOwner("relay-e2e")); err != nil {
			return nil, fmt.Errorf("could not apply object #%d %T: %+v", i, obj, err)
		}
	}

	return objs, nil
}

func WaitForServicesToBeReady(ctx context.Context, cl client.Client, namespace string) error {
	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		eps := &corev1.EndpointsList{}
		if err := cl.List(ctx, eps, client.InNamespace(namespace)); err != nil {
			return retry.RetryPermanent(err)
		}

		if len(eps.Items) == 0 {
			return retry.RetryTransient(fmt.Errorf("waiting for endpoints"))
		}

		for _, ep := range eps.Items {
			log.Println("checking service", ep.Name)
			if len(ep.Subsets) == 0 {
				return retry.RetryTransient(fmt.Errorf("waiting for subsets"))
			}

			for _, subset := range ep.Subsets {
				if len(subset.Addresses) == 0 {
					return retry.RetryTransient(fmt.Errorf("waiting for pod assignment"))
				}
			}
		}

		return retry.RetryPermanent(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func WaitForObjectDeletion(ctx context.Context, cl client.Client, obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	return retry.Retry(ctx, 1*time.Second, func() *retry.RetryError {
		if err := cl.Get(ctx, key, obj); errors.IsNotFound(err) {
			return retry.RetryPermanent(nil)
		} else if err != nil {
			return retry.RetryPermanent(err)
		}

		return retry.RetryTransient(fmt.Errorf("waiting for deletion of %T %s", obj, key))
	})
}

func SetKubernetesEnvVar(target *[]corev1.EnvVar, name, value string) {
	for i, ev := range *target {
		if ev.Name == name {
			(*target)[i].Value = value
			return
		}
	}

	*target = append(*target, corev1.EnvVar{Name: name, Value: value})
}
