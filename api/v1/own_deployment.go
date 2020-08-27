package v1

import appsv1 "k8s.io/api/apps/v1"

type OwnDeployment struct {
	Spec appsv1.DeploymentSpec `json:"spec"`
}
