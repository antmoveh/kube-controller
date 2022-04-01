/*
Copyright 2022 antmoveh.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	studyv1beta1 "github/antmoveh/kube-develop-tools/api/v1beta1"
)

// WorldReconciler reconciles a World object
type WorldReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	wf = "world.finalizers"
)

//+kubebuilder:rbac:groups=study.example.cn,resources=worlds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=study.example.cn,resources=worlds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=study.example.cn,resources=worlds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the World object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *WorldReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// your logic here
	wl := new(studyv1beta1.World)
	if err := r.Client.Get(ctx, req.NamespacedName, wl); err != nil {
		if !apierrs.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if wl.ObjectMeta.DeletionTimestamp == nil {
		if !containsString(wl.Finalizers, wf) {
			lv2 := wl.DeepCopy()
			lv2.Finalizers = append(lv2.Finalizers, wf)
			patch := client.MergeFrom(wl)
			if err := r.Patch(ctx, lv2, patch); err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("add finalizer")
			return ctrl.Result{Requeue: true}, nil
		}

		if wl.Status.War == "" {
			wl.Status.War = rand.String(8)
			wl.Status.SyncTime = metav1.Now()

			err := r.Client.Status().Update(ctx, wl)
			if err != nil {
				logger.Error(err, "update status failed")
			}
			return ctrl.Result{}, nil
		}
	}

	if !containsString(wl.Finalizers, wf) {
		return ctrl.Result{}, nil
	}

	lv2 := wl.DeepCopy()
	lv2.Finalizers = sliceRemoveString(lv2.Finalizers, wf)
	patch := client.MergeFrom(wl)
	if err := r.Patch(ctx, lv2, patch); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorldReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&studyv1beta1.World{}).
		Complete(r)
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func sliceRemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
