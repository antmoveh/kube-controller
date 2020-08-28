package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServicePort struct {
	Name       string             `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Protocol   string             `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`
	Port       int32              `json:"port" protobuf:"varint,3,opt,name=port"`
	TargetPort intstr.IntOrString `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
	NodePort   int32              `json:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
}

type OwnService struct {
	Ports     []corev1.ServicePort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,1,rep,name=ports"`
	ClusterIP string               `json:"clusterIp,omitempty" protobuf:"bytes,3,opt,name=clusterIp"`
}

type ServicePortStatus struct {
	corev1.ServicePort `json:"servicePort,omitempty"`
	Health             bool `json:"health,omitempty"`
}
