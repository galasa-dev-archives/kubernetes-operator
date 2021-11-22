package api

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Api) getDeployment() *appsv1.Deployment {
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
					InitContainers: []corev1.Container{
						{
							Name:            "init-chown-data",
							Image:           c.BusyboxImage,
							ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
							Command: []string{
								"chown", "-R", "1000", "/data",
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "data",
									MountPath: "/data",
									SubPath:   "",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            c.Name,
							Image:           c.Image,
							ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
							Command:         []string{"java"},
							Args: []string{
								"-jar",
								"boot.jar",
								"--obr",
								"file:galasa.obr",
								"--api",
								"--bootstrap",
								"file:/bootstrap.properties",
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
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 60,
								PeriodSeconds:       60,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(8080),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 3,
								PeriodSeconds:       1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "health",
										Port: intstr.FromInt(8080),
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "bootstrap",
									MountPath: "/bootstrap.properties",
									SubPath:   "bootstrap.properties",
								},
								{
									Name:      "testcatalog",
									MountPath: "/galasa/load/dev.galasa.testcatalog.cfg",
									SubPath:   "dev.galasa.testcatalog.cfg",
								},
								{
									Name:      "data",
									MountPath: "/galasa/testcatalog",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "bootstrap",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "bootstrap-file",
									},
								},
							},
						},
						{
							Name: "testcatalog",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "testcatalog-file",
									},
								},
							},
						},
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: c.Name + "-pvc",
								},
							},
						},
					},
				},
			},
		},
	}
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "Deployment"))
	return s
}
