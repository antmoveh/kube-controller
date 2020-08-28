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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	customv1 "kunit/api/v1"
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

func (r *UnitReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("unit", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *UnitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customv1.Unit{}).
		Complete(r)
}
