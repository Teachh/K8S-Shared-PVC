/*
Copyright 2024.

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

package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdv1 "github.com/Teachh/K8S-Shared-PVC/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// SharedPVCReconciler reconciles a SharedPVC object
type SharedPVCReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=crd.hector.dev,resources=sharedpvcs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.hector.dev,resources=sharedpvcs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=crd.hector.dev,resources=sharedpvcs/finalizers,verbs=update
// +kubebuilder:rbac:groups=crd.hector.dev,resources=sharedpvcs/events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=crd.hector.dev,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.hector.dev,resources=pods,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *SharedPVCReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	sharedpvc := &crdv1.SharedPVC{}
	// Get the SharedPVC and check the errors
	if err := r.Get(ctx, req.NamespacedName, sharedpvc); err != nil {
		l.Error(err, "Failed to get SharedPVC")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	l.Info(fmt.Sprintf("I've found: %s", sharedpvc.Name))
	// Modify in the future to not allow creating the manifest
	l.Info("Finding PVC")
	if err := r.checkIfPVCExists(ctx, sharedpvc); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	l.Info("Checking if the PVC can be mounted")
	if err := r.checkIfPVCCanBeMounted(ctx, sharedpvc); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	l.Info("Reconciliation Compelte!")

	return ctrl.Result{}, nil
}

func (r *SharedPVCReconciler) checkIfPVCExists(ctx context.Context, sharedpvc *crdv1.SharedPVC) error {
	l := log.FromContext(ctx)
	pvc := &corev1.PersistentVolumeClaim{}
	namespace := sharedpvc.Spec.NewPVC.OriginalNamespace
	name := sharedpvc.Spec.NewPVC.OriginalPVCName
	// Check if exists
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, pvc); err != nil {
		l.Error(err, fmt.Sprintf("Failed to find the PVC: %s", name))
		return err
	}
	// Check if the PV exists
	pv := &corev1.PersistentVolume{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: pvc.Spec.VolumeName}, pv); err != nil {
		l.Error(err, fmt.Sprintf("Failed to find the PV from the PVC: %s", pvc.Spec.VolumeName))
		return err
	}

	return nil
}

func (r *SharedPVCReconciler) checkIfPVCCanBeMounted(ctx context.Context, sharedpvc *crdv1.SharedPVC) error {
	l := log.FromContext(ctx)
	pvc := &corev1.PersistentVolumeClaim{}
	namespace := sharedpvc.Spec.NewPVC.OriginalNamespace
	name := sharedpvc.Spec.NewPVC.OriginalPVCName
	// Check if exists
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, pvc); err != nil {
		l.Error(err, fmt.Sprintf("Failed to find the PVC: %s", name))
		return err
	}

	podList := &corev1.PodList{}
	if err := r.Client.List(ctx, podList, &client.ListOptions{Namespace: pvc.Namespace}); err != nil {
		l.Error(err, "Error listing the pods")
		return err
	}

	// Iterate over each Pod to check if it is using the PVC
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			// Check last position
			if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvc.Name && !strings.Contains(string(pvc.Spec.AccessModes[len(pvc.Spec.AccessModes)-1]), "Only") {
				err := errors.New("PVC can not be used, check the Access mode and if it is already mounted")
				l.Error(err, "Error mounting the PVC")
			}
		}
	}
	return nil
}

// func (r *SharedPVCReconciler) createDeployment(ctx context.Context, sharedpvc *crdv1.SharedPVC) error {
// 	return nil
// }

// func (r *SharedPVCReconciler) createService(ctx context.Context, sharedpvc *crdv1.SharedPVC) error {
// 	return nil
// }

// SetupWithManager sets up the controller with the Manager.
func (r *SharedPVCReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Record Events from my controller
	r.recorder = mgr.GetEventRecorderFor("sharedpvc-controller")
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.SharedPVC{}).
		Complete(r)
}
