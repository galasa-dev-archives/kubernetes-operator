package metrics

import (
	"context"
	"reflect"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Metrics struct {
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

func New(metricsCrd *v2alpha1.GalasaMetricsComponent, k galasaecosystem.Interface) *Metrics {
	return &Metrics{
		Ecosystemclient: k,
		Namespace:       metricsCrd.Namespace,
		Name:            metricsCrd.Name,
		Image:           metricsCrd.Spec.Image,
		Replicas:        metricsCrd.Spec.Replicas,
		ImagePullPolicy: metricsCrd.Spec.ImagePullPolicy,
		NodeSelector:    metricsCrd.Spec.NodeSelector,
		Owner:           metricsCrd.OwnerReferences,
		Bootstrap:       metricsCrd.Spec.ComponentParms["bootstrap"],

		Status: v2alpha1.ComponentStatus{
			Ready: metricsCrd.Status.Ready,
		},
	}
}

func (c *Metrics) HasChanged(spec v2alpha1.ComponentSpec) bool {
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
func (c *Metrics) IsReady(ctx context.Context) bool {
	metrics, err := c.Ecosystemclient.GalasaV2alpha1().GalasaMetricsComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return metrics.Status.Ready
}
func (c *Metrics) GetObjects() []runtime.Object {
	var objects []runtime.Object

	objects = append(objects, c.getExposedService())
	objects = append(objects, c.getInternalService())
	objects = append(objects, c.getDeployment())

	return objects
}
