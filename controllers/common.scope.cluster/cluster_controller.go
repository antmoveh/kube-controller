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

package commonscopecluster

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	commonscopeclusterv1beta1 "github/antmoveh/kube-develop-tools/apis/common.scope.cluster/v1beta1"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=common.scope.cluster,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=common.scope.cluster,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=common.scope.cluster,resources=clusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	cu := new(commonscopeclusterv1beta1.Cluster)
	if err := r.Client.Get(ctx, req.NamespacedName, cu); err != nil {
		if !apierrs.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	r.Recorder.Event(cu, corev1.EventTypeNormal, "UpdateCluster", fmt.Sprintf("update cluster status %s", time.Now().Format("2006-01-02T15:04:05.000Z")))
	cu.Status.Cluster = rand.String(5)
	err := r.Client.Status().Update(ctx, cu)
	if err != nil {
		// log
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 创建联合索引，这样可以通过索引获取符合该索引条件的pod列表
	ctx := context.Background()
	err := mgr.GetFieldIndexer().IndexField(ctx, &corev1.Pod{}, "combinedIndex", func(object client.Object) []string {
		//combinedIndex := fmt.Sprintf("%s-%s", object.(*corev1.Pod).Spec.SchedulerName, object.(*corev1.Pod).Spec.NodeName)
		combinedIndex := fmt.Sprintf("%s-%s", object.(*corev1.Pod).Spec.SchedulerName, "")
		return []string{combinedIndex}
	})

	if err != nil {
		return err
	}

	pred := predicate.Funcs{
		CreateFunc: func(event.CreateEvent) bool { return true },
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc:  func(event.UpdateEvent) bool { return true },
		GenericFunc: func(event.GenericEvent) bool { return true },
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(pred).
		For(&commonscopeclusterv1beta1.Cluster{}, nodePredicateFn).
		WithOptions(controller.Options{
			RateLimiter: workqueue.NewItemFastSlowRateLimiter(10*time.Second, 60*time.Second, 5),
		}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}, podPredicateFn()).
		Complete(r)
}

var nodePredicateFn = builder.WithPredicates(
	predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			o := e.Object.(*commonscopeclusterv1beta1.Cluster)
			if o != nil {
				if len(o.Spec.ClusterName) > 0 {
					return true
				}
				return false
			}
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			o := e.Object.(*commonscopeclusterv1beta1.Cluster)
			if o != nil {
				if len(o.Spec.ClusterName) > 0 {
					return true
				}
				return false
			}
			return false
		},
		UpdateFunc:  func(event.UpdateEvent) bool { return false },
		GenericFunc: func(event.GenericEvent) bool { return false },
	})

func podPredicateFn() builder.Predicates {
	return builder.WithPredicates(predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			po := e.Object.(*corev1.Pod)
			if len(po.Name) > 5 {
				return true
			}
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			po := e.ObjectNew.(*corev1.Pod)
			if len(po.Name) > 5 {
				return true
			}
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			po := e.Object.(*corev1.Pod)
			if len(po.Name) > 5 {
				return true
			}
			return false
		},
		GenericFunc: func(event.GenericEvent) bool {
			return false
		},
	})
}

func (r *ClusterReconciler) getPod(ctx context.Context) error {

	podList := &corev1.PodList{}
	err := r.Client.List(ctx, podList, client.MatchingFields{"combinedIndex": fmt.Sprintf("%s-%s", "default-scheduler", "")})
	if err != nil {
		return err
	}
	for _, p := range podList.Items {
		fmt.Println("%s/%s", p.Namespace, p.Name)
	}
	return nil
}
