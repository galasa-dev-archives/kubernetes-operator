/*
 * Copyright contributors to the Galasa Project
 */
package toolbox

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Toolbox) getSimbankDeployment() *appsv1.Deployment {
	labels := map[string]string{
		"app": c.Name + "-simbank",
	}
	s := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name + "-simbank",
			Namespace:       c.Namespace,
			Labels:          labels,
			OwnerReferences: c.Owner,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: c.SimbankSpec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   c.Name + "-simbank",
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: c.SimbankSpec.NodeSelector,
					Containers: []corev1.Container{
						{
							Name:            "metrics",
							Image:           c.SimbankSpec.Image,
							ImagePullPolicy: corev1.PullPolicy(c.SimbankSpec.ImagePullPolicy),
							Command: []string{
								"java",
							},
							Args: []string{
								"-jar",
								"simplatform.jar",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "telnet",
									ContainerPort: 2023,
								},
								{
									Name:          "webservice",
									ContainerPort: 2080,
								},
								{
									Name:          "database",
									ContainerPort: 2027,
								},
								{
									Name:          "zosmf",
									ContainerPort: 2040,
								},
							},
						},
					},
					ServiceAccountName: "galasa-ecosystem-operator",
				},
			},
		},
	}
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "Deployment"))
	return s
}
