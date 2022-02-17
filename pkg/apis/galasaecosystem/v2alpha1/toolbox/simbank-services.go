package toolbox

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Toolbox) getSimbankDatabaseService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-simbank-database-external",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "simbank-database-external",
					TargetPort: intstr.FromInt(2027),
					Port:       2027,
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

func (c *Toolbox) getSimbankTelnetService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-simbank-telnet-external",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "simbank-telnet-external",
					TargetPort: intstr.FromInt(2023),
					Port:       2023,
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

func (c *Toolbox) getSimbankWebService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-simbank-webservice-external",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "simbank-webservice-external",
					TargetPort: intstr.FromInt(2080),
					Port:       2080,
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

func (c *Toolbox) getSimbankZosmfService() *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-simbank-zosmf-external",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(corev1.ServiceTypeNodePort),
			Ports: []corev1.ServicePort{
				{
					Name:       "simbank-zosmf-external",
					TargetPort: intstr.FromInt(2040),
					Port:       2040,
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
