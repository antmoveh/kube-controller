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
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	customv1 "kunit/api/v1"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OwnResource interface {
	// unit generate build-in resource
	MakeOwnResource(instance *customv1.Unit, logger logr.Logger, scheme *runtime.Scheme) (interface{}, error)
	// resource exits
	OwnResourceExist(instance *customv1.Unit, client client.Client, logger logr.Logger) (bool, interface{}, error)
	// update Unit status
	UpdateOwnResourceStatus(instance *customv1.Unit, client client.Client, logger logr.Logger) (*customv1.Unit, error)
	// create/update own build-in resource
	ApplyOwnResource(instance *customv1.Unit, client client.Client, logger logr.Logger, scheme *runtime.Scheme) error
}

// UnitReconciler reconciles a Unit object
type UnitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=custom.unit.crd.com,resources=units,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=custom.unit.crd.com,resources=units/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulSet,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployment,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=service,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoint,verbs=get
// +kubebuilder:rbac:groups=core,resources=persistentVolumeClaimStatus,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingress,verbs=get;update;patch;delete

func (r *UnitReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	//_ = context.Background()
	ctx := context.Background()
	_ = r.Log.WithValues("unit", req.NamespacedName)

	// your logic here
	defer func() {
		if rec := recover(); r != nil {
			switch x := rec.(type) {
			case error:
				r.Log.Error(x, "Reconcile error")
			}
		}
	}()

	// get unit object
	instance := &customv1.Unit{}

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
	}
	// 2 删除操作
	// 当资源对象Finalizer字段补位空是，delete操作会变成update操作，即为对象加上deletionTimeStamp时间戳
	// 当前时间在deletionTimestamp时间之后，且Finalizer已清空（视为清理后续任务处理完成）的情况下，gc会处理此对象

	myFinalizerName := "storage.finalizers.tutorial.kubebuilder.io"
	// DeletionTimestamp 时间戳为空，代表着当前对象不处于被删除的状态，为了开启Finalizer机制，先给他加上一段Finalizers,内容随机
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				r.Log.Error(err, "Add Finalizers error", instance.Namespace, instance.Name)
				return ctrl.Result{}, err
			}
		}
	} else {
		// DeletionTimestamp不为空，说明对象已经开始进入删除状态了。执行自己的删除步骤的后续逻辑，并清理掉自己的finalizer字段，等待自动gc
		if containsString(instance.ObjectMeta.Finalizers, myFinalizerName) {
			// 执行自定义删除逻辑
			if err := r.PreDelete(instance); err != nil {
				return ctrl.Result{}, err
			}
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, myFinalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// 3 创建或更新操作
	ownResources, err := r.getOwnResource(instance)
	if err != nil {
		msg := fmt.Sprintf("%s %s Reconciler.getOwnResource() function error", instance.Namespace, instance.Name)
		r.Log.Error(err, msg)
		return ctrl.Result{}, err
	}

	// 判断own resource是否存在
	success := true
	for _, ownResource := range ownResources {
		if err = ownResource.ApplyOwnResource(instance, r.Client, r.Log, r.Scheme); err != nil {
			success = false
		}
	}

	// 4 update Unit.status
	updateInstance := instance.DeepCopy()
	for _, ownResource := range ownResources {
		updateInstance, err = ownResource.UpdateOwnResourceStatus(updateInstance, r.Client, r.Log)
		if err != nil {
			success = false
		}
	}
	// apply update to apisServer if status changed
	if updateInstance != nil && !reflect.DeepEqual(updateInstance.Status, instance.Status) {
		if err := r.Status().Update(context.Background(), updateInstance); err != nil {
			r.Log.Error(err, "unable to update Unit status")
		}
	}
	// 5 记录结果
	if !success {
		msg := fmt.Sprintf("Reconciler Unit %s/%s failed", instance.Namespace, instance.Name)
		r.Log.Error(err, msg)
		return ctrl.Result{}, err
	} else {
		msg := fmt.Sprintf("Reconciler Unit %s/%s success", instance.Namespace, instance.Name)
		r.Log.Info(msg)
		return ctrl.Result{}, nil
	}
}

func (r *UnitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customv1.Unit{}).
		Complete(r)
}

func (r *UnitReconciler) PreDelete(instance *customv1.Unit) error {
	// 留空，自定义的pre delete逻辑需要，在这里实现
	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// 根据Unit.Spec生成其所有的own resource
func (r *UnitReconciler) getOwnResource(instance *customv1.Unit) ([]OwnResource, error) {
	var ownResources []OwnResource

	if instance.Spec.Category == "Deployment" {
		ownDeployment := customv1.OwnDeployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: instance.Spec.Replicas,
				Selector: instance.Spec.Selector,
				Template: instance.Spec.Template,
			},
		}
		ownDeployment.Spec.Template.Labels = instance.Spec.Selector.MatchLabels
		ownResources = append(ownResources, &ownDeployment)
	} else {
		ownStatefulSet := &customv1.OwnStatefulSet{
			Spec: appsv1.StatefulSetSpec{
				Replicas:    instance.Spec.Replicas,
				Selector:    instance.Spec.Selector,
				Template:    instance.Spec.Template,
				ServiceName: instance.Name,
			},
		}
		ownResources = append(ownResources, ownStatefulSet)
	}
	if instance.Spec.RelationResource.Service != nil {
		ownResources = append(ownResources, instance.Spec.RelationResource.Service)
	}
	if instance.Spec.RelationResource.Ingress != nil {
		ownResources = append(ownResources, instance.Spec.RelationResource.Ingress)
	}
	if instance.Spec.RelationResource.PVC != nil {
		ownResources = append(ownResources, instance.Spec.RelationResource.PVC)
	}
	return ownResources, nil
}
