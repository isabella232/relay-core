package obj

import (
	"context"

	relayv1beta1 "github.com/puppetlabs/relay-core/pkg/apis/relay.sh/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	TenantStatusReasonNamespaceReady = "NamespaceReady"
	TenantStatusReasonNamespaceError = "NamespaceError"

	TenantStatusReasonEventSinkMissing       = "EventSinkMissing"
	TenantStatusReasonEventSinkNotConfigured = "EventSinkNotConfigured"
	TenantStatusReasonEventSinkReady         = "EventSinkReady"

	TenantStatusReasonToolInjectionNotDefined = "ToolInjectionNotDefined"
	TenantStatusReasonToolInjectionError      = "ToolInjectionError"

	TenantStatusReasonReady = "Ready"
	TenantStatusReasonError = "Error"
)

var (
	TenantKind = relayv1beta1.SchemeGroupVersion.WithKind("Tenant")
)

type Tenant struct {
	Key    client.ObjectKey
	Object *relayv1beta1.Tenant
}

var _ Persister = &Tenant{}
var _ Finalizable = &Tenant{}
var _ Loader = &Tenant{}

func (t *Tenant) Persist(ctx context.Context, cl client.Client) error {
	return CreateOrUpdate(ctx, cl, t.Key, t.Object)
}

func (t *Tenant) PersistStatus(ctx context.Context, cl client.Client) error {
	return cl.Status().Update(ctx, t.Object)
}

func (t *Tenant) Finalizing() bool {
	return !t.Object.GetDeletionTimestamp().IsZero()
}

func (t *Tenant) AddFinalizer(ctx context.Context, name string) bool {
	return AddFinalizer(&t.Object.ObjectMeta, name)
}

func (t *Tenant) RemoveFinalizer(ctx context.Context, name string) bool {
	return RemoveFinalizer(&t.Object.ObjectMeta, name)
}

func (t *Tenant) Load(ctx context.Context, cl client.Client) (bool, error) {
	return GetIgnoreNotFound(ctx, cl, t.Key, t.Object)
}

func (t *Tenant) Managed() bool {
	return t.Object.Spec.NamespaceTemplate.Metadata.GetName() != ""
}

func (t *Tenant) Ready() bool {
	for _, cond := range t.Object.Status.Conditions {
		if cond.Type != relayv1beta1.TenantReady {
			continue
		}

		return cond.Status == corev1.ConditionTrue
	}

	return false
}

func NewTenant(key client.ObjectKey) *Tenant {
	return &Tenant{
		Key:    key,
		Object: &relayv1beta1.Tenant{},
	}
}

func ConfigureTenant(t *Tenant, td *TenantDepsResult, pvc *PersistentVolumeClaimResult) {
	// Set up our initial map from the existing data.
	conds := map[relayv1beta1.TenantConditionType]*relayv1beta1.Condition{
		relayv1beta1.TenantNamespaceReady:     &relayv1beta1.Condition{},
		relayv1beta1.TenantEventSinkReady:     &relayv1beta1.Condition{},
		relayv1beta1.TenantToolInjectionReady: &relayv1beta1.Condition{},
		relayv1beta1.TenantReady:              &relayv1beta1.Condition{},
	}

	for _, cond := range t.Object.Status.Conditions {
		*conds[cond.Type] = cond.Condition
	}

	// Update with dependency data.
	UpdateStatusConditionIfTransitioned(conds[relayv1beta1.TenantNamespaceReady], func() relayv1beta1.Condition {
		if td.Error != nil {
			return relayv1beta1.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  TenantStatusReasonNamespaceError,
				Message: td.Error.Error(),
			}
		} else if td.TenantDeps != nil {
			return relayv1beta1.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  TenantStatusReasonNamespaceReady,
				Message: "The tenant namespace is ready.",
			}
		}

		return relayv1beta1.Condition{
			Status: corev1.ConditionUnknown,
		}
	})

	UpdateStatusConditionIfTransitioned(conds[relayv1beta1.TenantEventSinkReady], func() relayv1beta1.Condition {
		if td.TenantDeps != nil {
			if sink := td.TenantDeps.APITriggerEventSink; sink != nil {
				if sink.URL() == "" {
					return relayv1beta1.Condition{
						Status:  corev1.ConditionFalse,
						Reason:  TenantStatusReasonEventSinkNotConfigured,
						Message: "The API trigger event sink is missing an endpoint URL.",
					}
				} else if _, ok := sink.Token(); !ok {
					return relayv1beta1.Condition{
						Status:  corev1.ConditionFalse,
						Reason:  TenantStatusReasonEventSinkNotConfigured,
						Message: "The API trigger event sink is missing a token.",
					}
				}

				return relayv1beta1.Condition{
					Status:  corev1.ConditionTrue,
					Reason:  TenantStatusReasonEventSinkReady,
					Message: "The event sink is ready.",
				}
			}

			// This shouldn't block people who want to use these APIs without
			// WebhookTriggers.
			return relayv1beta1.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  TenantStatusReasonEventSinkMissing,
				Message: "The tenant does not have an event sink defined.",
			}
		}

		return relayv1beta1.Condition{
			Status: corev1.ConditionUnknown,
		}
	})

	UpdateStatusConditionIfTransitioned(conds[relayv1beta1.TenantToolInjectionReady], func() relayv1beta1.Condition {
		if td.TenantDeps != nil {
			if vc := td.TenantDeps.ToolInjection.VolumeClaimTemplate; vc != nil {
				if pvc != nil {
					if pvc.Error != nil {
						return relayv1beta1.Condition{
							Status:  corev1.ConditionFalse,
							Reason:  TenantStatusReasonToolInjectionError,
							Message: pvc.Error.Error(),
						}
					} else if pvc.PersistentVolumeClaim != nil &&
						pvc.PersistentVolumeClaim.Object.Status.Phase == corev1.ClaimBound {
						return relayv1beta1.Condition{
							Status: corev1.ConditionTrue,
						}
					}
				}

				return relayv1beta1.Condition{
					Status: corev1.ConditionUnknown,
				}
			}

			return relayv1beta1.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  TenantStatusReasonToolInjectionNotDefined,
				Message: "The tenant tool injection in not defined.",
			}
		}

		return relayv1beta1.Condition{
			Status: corev1.ConditionUnknown,
		}
	})

	UpdateStatusConditionIfTransitioned(conds[relayv1beta1.TenantReady], func() relayv1beta1.Condition {
		switch AggregateStatusConditions(*conds[relayv1beta1.TenantNamespaceReady], *conds[relayv1beta1.TenantEventSinkReady], *conds[relayv1beta1.TenantToolInjectionReady]) {
		case corev1.ConditionTrue:
			return relayv1beta1.Condition{
				Status:  corev1.ConditionTrue,
				Reason:  TenantStatusReasonReady,
				Message: "The tenant is configured.",
			}
		case corev1.ConditionFalse:
			return relayv1beta1.Condition{
				Status:  corev1.ConditionFalse,
				Reason:  TenantStatusReasonError,
				Message: "One or more tenant components failed.",
			}
		}

		return relayv1beta1.Condition{
			Status: corev1.ConditionUnknown,
		}
	})

	t.Object.Status = relayv1beta1.TenantStatus{
		ObservedGeneration: t.Object.GetGeneration(),
		Namespace:          td.TenantDeps.Namespace.Name,
		Conditions: []relayv1beta1.TenantCondition{
			{
				Condition: *conds[relayv1beta1.TenantNamespaceReady],
				Type:      relayv1beta1.TenantNamespaceReady,
			},
			{
				Condition: *conds[relayv1beta1.TenantEventSinkReady],
				Type:      relayv1beta1.TenantEventSinkReady,
			},
			{
				Condition: *conds[relayv1beta1.TenantToolInjectionReady],
				Type:      relayv1beta1.TenantToolInjectionReady,
			},
			{
				Condition: *conds[relayv1beta1.TenantReady],
				Type:      relayv1beta1.TenantReady,
			},
		},
	}
}
