/*
 * Copyright contributors to the Galasa Project
 */
package resmon

import (
	"context"
	"reflect"

	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resmon struct {
	Ecosystemclient galasaecosystem.Interface
	Name            string
	Namespace       string

	Image           string
	Replicas        *int32
	ImagePullPolicy string
	NodeSelector    map[string]string
	Owner           []v1.OwnerReference
	Bootstrap       string
	Status          v2alpha1.ComponentStatus
}

func New(resmonCrd *v2alpha1.GalasaResmonComponent, k galasaecosystem.Interface) *Resmon {
	t := true
	return &Resmon{
		Ecosystemclient: k,
		Namespace:       resmonCrd.Namespace,
		Name:            resmonCrd.Name,
		Image:           resmonCrd.Spec.Image,
		Replicas:        resmonCrd.Spec.Replicas,
		ImagePullPolicy: resmonCrd.Spec.ImagePullPolicy,
		NodeSelector:    resmonCrd.Spec.NodeSelector,
		Owner: []v1.OwnerReference{
			{
				APIVersion:         "galasa.dev/v2alpha1",
				Kind:               "GalasaResmonComponent",
				Name:               resmonCrd.Name,
				UID:                resmonCrd.GetUID(),
				Controller:         &t,
				BlockOwnerDeletion: &t,
			},
		},
		Bootstrap: resmonCrd.Spec.ComponentParms["bootstrap"],

		Status: v2alpha1.ComponentStatus{
			Ready: resmonCrd.Status.Ready,
		},
	}
}

func (c *Resmon) HasChanged(spec v2alpha1.ComponentSpec) bool {
	if c.Image != spec.Image {
		return true
	}
	if c.ImagePullPolicy != spec.ImagePullPolicy {
		return true
	}
	if reflect.DeepEqual(c.NodeSelector, spec.NodeSelector) {
		return true
	}
	return false
}
func (c *Resmon) IsReady(ctx context.Context) bool {
	resmon, err := c.Ecosystemclient.GalasaV2alpha1().GalasaResmonComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return resmon.Status.Ready
}
func (c *Resmon) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getExposedService())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getDeployment())

	return objects
}
