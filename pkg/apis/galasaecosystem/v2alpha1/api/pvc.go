package api

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Api) getPersistentVolumeClaim() *corev1.PersistentVolumeClaim {
	p := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:       c.Name + "-pvc",
			Namespace:  c.Namespace,
			Finalizers: []string{"kubernetes.io/pvc-protection"},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			StorageClassName: &c.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(c.Storage),
				},
			},
		},
	}

	p.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "PersistentVolumeClaim"))
	return p
}
