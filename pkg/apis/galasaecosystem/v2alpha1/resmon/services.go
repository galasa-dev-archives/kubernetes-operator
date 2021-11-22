package resmon

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Resmon) getExposedService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-external-service",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "metrics",
					TargetPort: intstr.FromInt(9010),
					Port:       9010,
				},
				{
					Name:       "health",
					TargetPort: intstr.FromInt(9011),
					Port:       9011,
				},
			},
			Selector: map[string]string{
				"app": c.Name,
			},
		},
	}
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "Service"))
	return s
}

func (c *Resmon) getInternalService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-internal-service",
			Namespace: c.Namespace,
			Labels: map[string]string{
				"app": c.Name,
			},
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "metrics",
					TargetPort: intstr.FromInt(9010),
					Port:       9010,
				},
				{
					Name:       "health",
					TargetPort: intstr.FromInt(9011),
					Port:       9011,
				},
			},
			Selector: map[string]string{
				"app": c.Name,
			},
		},
	}
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "Service"))
	return s
}
