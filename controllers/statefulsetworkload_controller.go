/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cpv1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	corev1alpha2 "statefulsetworkload/api/v1alpha2"
)

const (
	oamReconcileWait = 30 * time.Second
)

//Reconcile error strings
const (
	errRenderWorkload   = "cannot render workload"
	errUpdateStatus     = "cannot apply status"
	errApplyStatefulSet = "cannot apply the statefulset"
	errGCStatefulSet    = "cannot clean up stale statefulsets"
)

// StatefulSetWorkloadReconciler reconciles a StatefulSetWorkload object
type StatefulSetWorkloadReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.oam.dev,resources=statefulsetworkloads,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.oam.dev,resources=statefulsetworkloads/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

func (r *StatefulSetWorkloadReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("statefulsetworkload", req.NamespacedName)
	log.Info("Reconcile container workload")

	var workload corev1alpha2.StatefulSetWorkload
	if err := r.Get(ctx, req.NamespacedName, &workload); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Container wokload is deleted")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Get the workload", "apiVersion", workload.APIVersion, "kind", workload.Kind)

	statefulset, err := r.renderStatefulSet(ctx, &workload)
	if err != nil {
		workload.Status.SetConditions(cpv1alpha1.ReconcileError(errors.Wrap(err, errRenderWorkload)))
		log.Error(err, "Failed to render a statefulset")
		return reconcile.Result{RequeueAfter: oamReconcileWait}, errors.Wrap(r.Status().Update(ctx, &workload),
			errUpdateStatus)
	}
	log.Info("Get a statefulset", "statefulset", statefulset.Spec.Template.Spec.Containers[0])

	log.Info("Successfully rendered a statefulset",
		"statefulset name", statefulset.Name,
		"statefulset Namespace", statefulset.Namespace,
		"number of containers", len(statefulset.Spec.Template.Spec.Containers),
		"first container image", statefulset.Spec.Template.Spec.Containers[0].Image)
	// server side apply, only the fields we set are touched
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner(workload.Name)}
	if err := r.Patch(ctx, statefulset, client.Apply, applyOpts...); err != nil {
		workload.Status.SetConditions(cpv1alpha1.ReconcileError(errors.Wrap(err, errApplyStatefulSet)))
		log.Error(err, "Failed to apply to a statefulset")
		return reconcile.Result{RequeueAfter: oamReconcileWait}, errors.Wrap(r.Status().Update(ctx, &workload),
			errUpdateStatus)
	}
	log.Info("Successfully applied a statefulset", "UID", statefulset.UID)

	if err := r.Status().Update(ctx, &workload); err != nil {
		return reconcile.Result{RequeueAfter: oamReconcileWait}, err
	}

	// create a service for the workload
	// TODO(rz): Use ingress trait instead

	// garbage collect the statefulsets that we created but not needed
	if err := r.cleanupResources(ctx, &workload, &statefulset.UID); err != nil {
		workload.Status.SetConditions(cpv1alpha1.ReconcileError(errors.Wrap(err, errGCStatefulSet)))
		log.Error(err, "Failed to clean up resources")
		return reconcile.Result{RequeueAfter: oamReconcileWait}, errors.Wrap(r.Status().Update(ctx, &workload),
			errUpdateStatus)
	}
	workload.Status.Resources = nil
	// record the new statefulset
	workload.Status.Resources = append(workload.Status.Resources, cpv1alpha1.TypedReference{
		APIVersion: statefulset.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		Kind:       statefulset.GetObjectKind().GroupVersionKind().Kind,
		Name:       statefulset.GetName(),
		UID:        statefulset.GetUID(),
	})
	// record the new service

	if err := r.Status().Update(ctx, &workload); err != nil {
		return reconcile.Result{RequeueAfter: oamReconcileWait}, err
	}

	workload.Status.SetConditions(cpv1alpha1.ReconcileSuccess())
	return ctrl.Result{}, errors.Wrap(r.Status().Update(ctx, &workload), errUpdateStatus)
}

func (r *StatefulSetWorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha2.StatefulSetWorkload{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
