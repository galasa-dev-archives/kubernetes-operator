package ras

import (
	"context"
	"reflect"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Ras struct {
	Ecosystemclient galasaecosystem.Interface
	Name            string
	Namespace       string

	Image            string
	Replicas         *int32
	ImagePullPolicy  string
	Storage          string
	StorageClassName string
	NodeSelector     map[string]string
	Owner            []v1.OwnerReference
	Status           v2alpha1.ComponentStatus
}

func New(rasCrd *v2alpha1.GalasaRasComponent, k galasaecosystem.Interface) *Ras {
	return &Ras{
		Ecosystemclient:  k,
		Namespace:        rasCrd.Namespace,
		Name:             rasCrd.Name,
		Image:            rasCrd.Spec.Image,
		Replicas:         rasCrd.Spec.Replicas,
		ImagePullPolicy:  rasCrd.Spec.ImagePullPolicy,
		Storage:          rasCrd.Spec.Storage,
		StorageClassName: rasCrd.Spec.StorageClassName,
		NodeSelector:     rasCrd.Spec.NodeSelector,
		Owner:            rasCrd.OwnerReferences,

		Status: v2alpha1.ComponentStatus{
			Ready: rasCrd.Status.Ready,
		},
	}
}

func (c *Ras) HasChanged(spec v2alpha1.ComponentSpec) bool {
	if c.Image != spec.Image {
		return true
	}
	if c.ImagePullPolicy != spec.ImagePullPolicy {
		return true
	}
	if c.Storage != spec.Storage {
		return true
	}
	if c.StorageClassName != spec.StorageClassName {
		return true
	}
	if reflect.DeepEqual(c.NodeSelector, spec.NodeSelector) {
		return true
	}
	return false
}
func (c *Ras) IsReady(ctx context.Context) bool {
	ras, err := c.Ecosystemclient.GalasaV2alpha1().GalasaRasComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return ras.Status.Ready
}
func (c *Ras) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getPersistentVolumeClaim())
	objects = append(objects, c.getExposedService())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getStatefulSet())

	return objects
}
