package ras

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *Ras) getStatefulSet() *appsv1.StatefulSet {
	labels := map[string]string{
		"app": c.Name,
	}
	s := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            c.Name,
			Namespace:       c.Namespace,
			Labels:          labels,
			OwnerReferences: c.Owner,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: c.Name + "-internal-service",
			Replicas:    c.Replicas,
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
							Name:            "couchdb",
							Image:           c.Image,
							ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy),
							Ports: []corev1.ContainerPort{
								{
									Name:          "couchdbport",
									ContainerPort: 5984,
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 60,
								PeriodSeconds:       60,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(5984),
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/opt/couchdb/data",
									Name:      "data-disk",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data-disk",
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
	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("apps/v1", "StatefulSet"))
	return s
}
