package toolbox

/*
 * Copyright contributors to the Galasa Project
 */

import (
	"context"

	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	galasaecosystem "github.com/galasa-dev/kubernetes-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Toolbox struct {
	Ecosystemclient galasaecosystem.Interface
	Name            string
	Namespace       string

	Simbank     bool
	SimbankSpec v2alpha1.ComponentSpec

	//TODO
	// Prometheus bool
	// PrometheusSpec v2alpha1.ComponentSpec

	// Grafana bool
	// GradanaSpec v2alpha1.ComponentSpec

	// Nexus bool
	// NexusSpec v2alpha1.ComponentSpec

	// Jenkins bool
	// JenkinsSpec v2alpha1.ComponentSpec

	ImagePullPolicy string
	NodeSelector    map[string]string
	Owner           []v1.OwnerReference
	Status          v2alpha1.ComponentStatus
}

func New(toolboxCrd *v2alpha1.GalasaToolboxComponent, k galasaecosystem.Interface) *Toolbox {
	t := true
	return &Toolbox{
		Ecosystemclient: k,
		Namespace:       toolboxCrd.Namespace,
		Name:            toolboxCrd.Name,
		Owner: []v1.OwnerReference{
			{
				APIVersion:         "galasa.dev/v2alpha1",
				Kind:               "GalasaToolboxComponent",
				Name:               toolboxCrd.Name,
				UID:                toolboxCrd.GetUID(),
				Controller:         &t,
				BlockOwnerDeletion: &t,
			},
		},

		Simbank:     toolboxCrd.Spec.Simbank,
		SimbankSpec: toolboxCrd.Spec.SimbankSpec,

		Status: v2alpha1.ComponentStatus{
			Ready: toolboxCrd.Status.Ready,
		},
	}
}

func (c *Toolbox) HasChanged(spec v2alpha1.ComponentSpec) bool {
	return false
}

func (c *Toolbox) IsReady(ctx context.Context) bool {
	toolbox, err := c.Ecosystemclient.GalasaV2alpha1().GalasaToolboxComponents(c.Namespace).Get(ctx, c.Name, v1.GetOptions{})
	if err != nil {
		return false
	}
	return toolbox.Status.Ready
}
func (c *Toolbox) GetObjects() []runtime.Object {
	var objects []runtime.Object

	if c.Simbank {
		// Simbank
		objects = append(objects, c.getSimbankDatabaseService())
		objects = append(objects, c.getSimbankTelnetService())
		objects = append(objects, c.getSimbankWebService())
		objects = append(objects, c.getSimbankZosmfService())
		objects = append(objects, c.getSimbankDeployment())
	}

	return objects
}
