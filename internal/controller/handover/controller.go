package handover

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	coordinationv1alpha1 "github.com/thetechnick/package-operator/apis/coordination/v1alpha1"
	"github.com/thetechnick/package-operator/internal/dynamicwatcher"
	internalprobe "github.com/thetechnick/package-operator/internal/probe"
)

type HandoverReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	DynamicClient   dynamic.Interface
	DiscoveryClient *discovery.DiscoveryClient

	dw *dynamicwatcher.DynamicWatcher
}

const (
	cacheFinalizer = "coordination.thetechnick.ninja/cache"
)

func (r *HandoverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.dw = dynamicwatcher.New(r.Log, r.Scheme, r.RESTMapper(), r.DynamicClient)

	return ctrl.NewControllerManagedBy(mgr).
		For(&coordinationv1alpha1.Handover{}).
		Watches(r.dw, &handler.EnqueueRequestForOwner{
			OwnerType:    &coordinationv1alpha1.Handover{},
			IsController: false,
		}).
		Complete(r)
}

const (
	maxParrallel = 1
)

func (r *HandoverReconciler) Reconcile(
	ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Handover", req.NamespacedName.String())

	handover := &coordinationv1alpha1.Handover{}
	if err := r.Get(ctx, req.NamespacedName, handover); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !handover.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.handleDeletion(ctx, handover)
	}

	// Add finalizers
	if !controllerutil.ContainsFinalizer(
		handover, cacheFinalizer) {
		controllerutil.AddFinalizer(handover, cacheFinalizer)
		if err := r.Update(ctx, handover); err != nil {
			return ctrl.Result{}, fmt.Errorf("adding finalizer: %w", err)
		}
	}

	// Ensure watch
	gvk := schema.GroupVersionKind{
		Group:   handover.Spec.Target.Group,
		Version: handover.Spec.Target.Version,
		Kind:    handover.Spec.Target.Kind,
	}
	objType := &unstructured.Unstructured{}
	objType.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   handover.Spec.Target.Group,
		Version: handover.Spec.Target.Version,
		Kind:    handover.Spec.Target.Kind,
	})
	if err := r.dw.Watch(handover, objType); err != nil {
		return ctrl.Result{}, fmt.Errorf("watching %s: %w", gvk, err)
	}

	// Parse Probes
	probe := internalprobe.Parse(handover.Spec.Probes)

	// Handle processing objects
	if err := r.handleProcessing(ctx, log, handover, objType, probe); err != nil {
		return ctrl.Result{}, err
	}

	// List all objects
	// select all objects with new or old label value
	requirement, err := labels.NewRequirement(
		handover.Spec.Strategy.Relabel.Label,
		selection.In,
		[]string{
			handover.Spec.Strategy.Relabel.NewValue,
			handover.Spec.Strategy.Relabel.OldValue,
		},
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("building selector: %w", err)
	}
	selector := labels.NewSelector().Add(*requirement)

	objList := &unstructured.UnstructuredList{}
	objList.SetGroupVersionKind(gvk)
	objList.SetKind(gvk.Kind + "List")
	if err := r.List(
		ctx, objList,
		client.InNamespace(handover.Namespace),
		client.MatchingLabelsSelector{
			Selector: selector,
		},
	); err != nil {
		return ctrl.Result{}, fmt.Errorf("listing %s: %w", gvk, err)
	}

	// Check state
	var unavailable int
	for _, obj := range objList.Items {
		if success, _ := probe.Probe(&obj); !success {
			unavailable++
		}
	}

	// split into old and new
	groups := groupByLabelValues(
		objList.Items, handover.Spec.Strategy.Relabel.Label,
		handover.Spec.Strategy.Relabel.NewValue,
		handover.Spec.Strategy.Relabel.OldValue,
	)
	newObjs := groups[0]
	oldObjs := groups[1]

	if unavailable > 0 {
		// Don't process anymore when something is unavailable.
		for _, obj := range oldObjs {
			if len(handover.Status.Processing) < maxParrallel {
				break
			}

			handover.Status.Processing = append(
				handover.Status.Processing,
				coordinationv1alpha1.HandoverRef{
					UID:  obj.GetUID(),
					Name: obj.GetName(),
				})
		}
	}

	// report counts
	handover.Status.Stats.Found = int32(len(objList.Items))
	handover.Status.Stats.Updated = int32(len(newObjs))
	handover.Status.Stats.Available = handover.Status.Stats.Found - int32(unavailable)

	if handover.Status.Stats.Found == handover.Status.Stats.Updated {
		meta.SetStatusCondition(&handover.Status.Conditions, metav1.Condition{
			Type:               coordinationv1alpha1.HandoverCompleted,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: handover.Generation,
			Reason:             "Complete",
			Message:            "All found objects have been re-labeled.",
		})
	} else {
		meta.SetStatusCondition(&handover.Status.Conditions, metav1.Condition{
			Type:               coordinationv1alpha1.HandoverCompleted,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: handover.Generation,
			Reason:             "Incomplete",
			Message:            "Some found objects need to be re-labeled.",
		})
	}

	if err := r.Status().Update(ctx, handover); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating Handover status: %w", err)
	}
	return ctrl.Result{}, nil
}

func (r *HandoverReconciler) handleProcessing(
	ctx context.Context,
	log logr.Logger,
	handover *coordinationv1alpha1.Handover,
	objType *unstructured.Unstructured,
	probe internalprobe.Interface,
) error {
	var stillProcessing []coordinationv1alpha1.HandoverRef
	for _, processing := range handover.Status.Processing {
		processingObj := objType.DeepCopy()
		err := r.Client.Get(ctx, types.NamespacedName{
			Name:      processing.Name,
			Namespace: handover.Namespace,
		}, processingObj)
		if err != nil {
			return fmt.Errorf("getting object in process queue: %w", err)
		}

		labels := processingObj.GetLabels()
		if labels == nil ||
			labels[handover.Spec.Strategy.Relabel.Label] != handover.Spec.Strategy.Relabel.NewValue {
			labels[handover.Spec.Strategy.Relabel.Label] = handover.Spec.Strategy.Relabel.NewValue
			processingObj.SetLabels(labels)
			if err := r.Update(ctx, processingObj); err != nil {
				return fmt.Errorf("updating object in process queue: %w", err)
			}
		}

		if success, message := probe.Probe(processingObj); !success {
			log.Info("waiting to be ready", "objName", processing.Name, "failure", message)
			stillProcessing = append(stillProcessing, processing)
		}
	}
	handover.Status.Processing = stillProcessing
	return nil
}

func groupByLabelValues(in []unstructured.Unstructured, labelKey string, values ...string) [][]unstructured.Unstructured {
	out := make([][]unstructured.Unstructured, len(values))
	for _, obj := range in {
		if obj.GetLabels() == nil {
			continue
		}
		if len(obj.GetLabels()[labelKey]) == 0 {
			continue
		}

		for i, v := range values {
			if obj.GetLabels()[labelKey] == v {
				out[i] = append(out[i], obj)
			}
		}
	}
	return out
}

// handle deletion of the PackageSet
func (r *HandoverReconciler) handleDeletion(
	ctx context.Context,
	handover *coordinationv1alpha1.Handover,
) error {
	if controllerutil.ContainsFinalizer(
		handover, cacheFinalizer) {
		controllerutil.RemoveFinalizer(
			handover, cacheFinalizer)

		if err := r.Update(ctx, handover); err != nil {
			return fmt.Errorf("removing finalizer: %w", err)
		}
	}

	if err := r.dw.Free(handover); err != nil {
		return fmt.Errorf("free cache: %w", err)
	}
	return nil
}
