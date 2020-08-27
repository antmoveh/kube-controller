package v1

import appsv1 "k8s.io/api/apps/v1"

type OwnStatefulSet struct {
	Spec appsv1.StatefulSetSpec
}
