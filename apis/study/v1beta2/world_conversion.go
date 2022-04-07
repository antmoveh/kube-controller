package v1beta2

import (
	"github/antmoveh/kube-develop-tools/apis/study/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *World) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.World)
	dst.Spec.World = src.Spec.Earth
	return nil
}

func (dst *World) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.World)
	dst.Spec.Earth = src.Spec.World
	return nil
}
