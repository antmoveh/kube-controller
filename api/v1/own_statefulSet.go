package v1

import (
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OwnStatefulSet struct {
	Spec appsv1.StatefulSetSpec
}

func (ownStatefulset *OwnStatefulSet) MakeOwnResource(instance *Unit, logger logr.Logger, scheme *runtime.Scheme) (interface{}, error) {

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: instance.Name, Namespace: instance.Namespace, Labels: instance.Labels},
		Spec:       ownStatefulset.Spec,
	}

	customizeEnvs := []corev1.EnvVar{
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"},
			},
		},
		{
			Name:  "APPNAME",
			Value: instance.Name,
		},
	}

	var specEnvs []corev1.EnvVar
	templateEnvs := sts.Spec.Template.Spec.Containers[0].Env
	for index := range templateEnvs {
		if templateEnvs[index].Name != "POD_NAME" && templateEnvs[index].Name != "APPNAME" {
			specEnvs = append(specEnvs, templateEnvs[index])
		}
	}

	sts.Spec.Template.Spec.Containers[0].Env = append(specEnvs, customizeEnvs...)

	if err := controllerutil.SetControllerReference(instance, sts, scheme); err != nil {
		msg := fmt.Sprintf("set controllerReference for StatefulSet %s%s failed", instance.Namespace, instance.Name)
		logger.Error(err, msg)
		return nil, err
	}
	return sts, nil
}
