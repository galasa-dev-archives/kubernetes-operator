/*
 * Copyright contributors to the Galasa Project
 */
package cps

import (
	"context"
	"reflect"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Cps struct {
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

func New(cpsCrd *v2alpha1.GalasaCpsComponent, k galasaecosystem.Interface) *Cps {
	return &Cps{
		Ecosystemclient:  k,
		Namespace:        cpsCrd.Namespace,
		Name:             cpsCrd.Name,
		Image:            cpsCrd.Spec.Image,
		Replicas:         cpsCrd.Spec.Replicas,
		ImagePullPolicy:  cpsCrd.Spec.ImagePullPolicy,
		Storage:          cpsCrd.Spec.Storage,
		StorageClassName: cpsCrd.Spec.StorageClassName,
		NodeSelector:     cpsCrd.Spec.NodeSelector,
		Owner:            cpsCrd.OwnerReferences,

		Status: v2alpha1.ComponentStatus{
			Ready: cpsCrd.Status.Ready,
		},
	}
}

func (c *Cps) HasChanged(spec v2alpha1.ComponentSpec) bool {
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
func (c *Cps) IsReady(ctx context.Context) bool {
	cps, err := c.Ecosystemclient.GalasaV2alpha1().GalasaCpsComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return cps.Status.Ready
}
func (c *Cps) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getPersistentVolumeClaim())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getExposedService())
	objects = append(objects, c.getStatefulSet())

	return objects
}
