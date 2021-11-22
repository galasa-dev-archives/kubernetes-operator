package api

import (
	"context"
	"reflect"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Api struct {
	Ecosystemclient galasaecosystem.Interface
	Name            string
	Namespace       string

	BusyboxImage     string
	Image            string
	Replicas         *int32
	ImagePullPolicy  string
	Storage          string
	StorageClassName string
	NodeSelector     map[string]string
	Owner            []v1.OwnerReference
	Status           v2alpha1.ComponentStatus

	CPSUri string
}

func New(apiCrd *v2alpha1.GalasaApiComponent, k galasaecosystem.Interface) *Api {
	return &Api{
		Ecosystemclient:  k,
		Namespace:        apiCrd.Namespace,
		Name:             apiCrd.Name,
		BusyboxImage:     apiCrd.Spec.ComponentParms["busyboxImage"],
		CPSUri:           apiCrd.Spec.ComponentParms["cpsuri"],
		Image:            apiCrd.Spec.Image,
		Replicas:         apiCrd.Spec.Replicas,
		ImagePullPolicy:  apiCrd.Spec.ImagePullPolicy,
		Storage:          apiCrd.Spec.Storage,
		StorageClassName: apiCrd.Spec.StorageClassName,
		NodeSelector:     apiCrd.Spec.NodeSelector,
		Owner:            apiCrd.OwnerReferences,

		Status: v2alpha1.ComponentStatus{
			Ready: apiCrd.Status.Ready,
		},
	}
}

func (c *Api) HasChanged(spec v2alpha1.ComponentSpec) bool {
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

func (c *Api) IsReady(ctx context.Context) bool {
	api, err := c.Ecosystemclient.GalasaV2alpha1().GalasaApiComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return api.Status.Ready
}
func (c *Api) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getPersistentVolumeClaim())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getExposedService())
	objects = append(objects, c.getTestCatalog())
	objects = append(objects, c.getBootstrap())
	objects = append(objects, c.getDeployment())

	return objects
}
