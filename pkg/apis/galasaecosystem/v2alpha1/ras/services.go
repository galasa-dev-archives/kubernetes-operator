package ras

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Ras) getExposedService() *corev1.Service {
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
					Name:       "couchdbport",
					TargetPort: intstr.FromInt(5984),
					Port:       5984,
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

func (c *Ras) getInternalService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-internal-service",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
			Labels: map[string]string{
				"app": c.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "couchdbport",
					TargetPort: intstr.FromInt(5984),
					Port:       5984,
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
