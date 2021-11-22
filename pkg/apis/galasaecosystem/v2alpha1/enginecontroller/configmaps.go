package enginecontroller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *EngineController) getControllerConfig() *corev1.ConfigMap {
	config := map[string]string{
		"bootstrap":    c.Bootstrap,
		"max_engines":  "10",
		"engine_label": "k8s-standard-engine",
		"node_arch":    "amd64",
		"engine_image": c.Image,
	}

	s := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "config",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Data: config,
	}

	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "ConfigMap"))
	return s
}
