package metrics

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Metrics) getDeployment() *appsv1.Deployment {
	labels := map[string]string{
		"app": c.Name,
	}
	s := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name,
			Namespace:       c.Namespace,
			Labels:          labels,
			OwnerReferences: c.Owner,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: c.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   c.Name,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: c.NodeSelector,
					Containers: []corev1.Container{
						{
							Name:            "metrics",
							Image:           c.Image,
							ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
							Command: []string{
								"java",
							},
							Args: []string{
								"-jar",
								"boot.jar",
								"--obr",
								"file:galasa.obr",
								"--trace",
								"--metricserver",
								"--bootstrap",
								"$(BOOTSTRAP_URI)",
								"--trace",
							},
							Env: []corev1.EnvVar{
								{
									Name: "BOOTSTRAP_URI",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "config",
											},
											Key: "bootstrap",
										},
									},
								},
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 9010,
								},
								{
									Name:          "health",
									ContainerPort: 9011,
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 60,
								PeriodSeconds:       60,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(9011),
									},
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
