/*
 * Copyright contributors to the Galasa Project
 */
package enginecontroller

import (
	"context"
	"reflect"

	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type EngineController struct {
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

func New(engineControllerCrd *v2alpha1.GalasaEngineControllerComponent, k galasaecosystem.Interface) *EngineController {
	t := true
	return &EngineController{
		Ecosystemclient: k,
		Namespace:       engineControllerCrd.Namespace,
		Name:            engineControllerCrd.Name,
		Image:           engineControllerCrd.Spec.Image,
		Replicas:        engineControllerCrd.Spec.Replicas,
		ImagePullPolicy: engineControllerCrd.Spec.ImagePullPolicy,
		NodeSelector:    engineControllerCrd.Spec.NodeSelector,
		Owner: []v1.OwnerReference{
			{
				APIVersion:         "galasa.dev/v2alpha1",
				Kind:               "GalasaEngineControllerComponent",
				Name:               engineControllerCrd.Name,
				UID:                engineControllerCrd.GetUID(),
				Controller:         &t,
				BlockOwnerDeletion: &t,
			},
		},
		Bootstrap: engineControllerCrd.Spec.ComponentParms["bootstrap"],

		Status: v2alpha1.ComponentStatus{
			Ready: engineControllerCrd.Status.Ready,
		},
	}
}

func (c *EngineController) HasChanged(spec v2alpha1.ComponentSpec) bool {
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
func (c *EngineController) IsReady(ctx context.Context) bool {
	ec, err := c.Ecosystemclient.GalasaV2alpha1().GalasaEngineControllerComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return ec.Status.Ready
}
func (c *EngineController) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getControllerConfig())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getDeployment())

	return objects
}
