package api

import (
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Api) getBootstrap() *corev1.ConfigMap {
	s := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "bootstrap-file",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Data: map[string]string{
			"bootstrap.properties": `framework.config.store=etcd:` + c.CPSUri + `
framework.extra.bundles=dev.galasa.cps.etcd,dev.galasa.ras.couchdb
			`,
		},
	}

	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "ConfigMap"))
	return s
}

func (c *Api) getTestCatalog() *corev1.ConfigMap {
	s := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "testcatalog-file",
			Namespace:       c.Namespace,
			OwnerReferences: c.Owner,
		},
		Data: map[string]string{
			"dev.galasa.testcatalog.cfg": "framework.testcatalog.directory=file:/galasa/testcatalog",
		},
	}

	s.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "ConfigMap"))
	return s
}
