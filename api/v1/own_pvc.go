package v1

import corev1 "k8s.io/api/core/v1"

type OwnPVC struct {
	Spec corev1.PersistentVolumeClaimSpec
}
