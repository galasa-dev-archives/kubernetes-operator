package cps

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Cps) getExposedService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-external-service",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
			Labels: map[string]string{
				"app": c.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "etcd-client",
					TargetPort: intstr.FromInt(2379),
					Port:       2379,
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

func (c *Cps) getInternalService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
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
					Name:       "etcd-server",
					TargetPort: intstr.FromInt(2379),
					Port:       2379,
				},
				{
					Name:       "etcd-client",
					TargetPort: intstr.FromInt(2380),
					Port:       2380,
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
