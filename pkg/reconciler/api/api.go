/*
 * Copyright contributors to the Galasa Project
 */
package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	ecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	galasaecosystemlisters "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/listers/galasaecosystem/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
)

type Reconciler struct {
	KubeClientSet kubernetes.Interface

	GalasaEcosystemClientSet galasaecosystem.Interface
	GalasaApiLister          galasaecosystemlisters.GalasaApiComponentLister
	GalasaCpsLister          galasaecosystemlisters.GalasaCpsComponentLister
}

func (c *Reconciler) ReconcileKind(ctx context.Context, p *v2alpha1.GalasaApiComponent) pkgreconciler.Event {
	logger := logging.FromContext(ctx)
	api := ecosystem.Api(p, c.GalasaEcosystemClientSet)
	objects := api.GetObjects()
	logger.Infof("Reconciling API")

	for _, obj := range objects {
		switch obj.GetObjectKind().GroupVersionKind() {
		case schema.FromAPIVersionAndKind("v1", "Service"):
			logger.Infof("Found service: %s", obj.(*corev1.Service).Name)
			service := obj.(*corev1.Service)
			s, _ := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, service.Name, v1.GetOptions{})
			if s.Name == "" {
				logger.Infof("Create service")
				_, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Create(ctx, service, v1.CreateOptions{})
				if err != nil {
					return controller.NewPermanentError(fmt.Errorf("failed to create service resources: %v", err))
				}
			} else {
				logger.Infof("Service pre-existing, please manually remove service %s to apply new changes: %v", service.Name, s.Name)
			}

		case schema.FromAPIVersionAndKind("apps/v1", "Deployment"):
			logger.Infof("Found Deployment: %s", obj.(*appsv1.Deployment).Name)
			d := obj.(*appsv1.Deployment)
			s, _ := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Get(ctx, d.Name, v1.GetOptions{})
			if s.Name == "" {
				logger.Infof("Create Deployment")
				_, err := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Create(ctx, d, v1.CreateOptions{})
				if err != nil {
					return controller.NewPermanentError(fmt.Errorf("failed to create Deployments resources: %v", err))
				}
			} else {
				logger.Infof("Updating Deployment with new configuration")
				_, err := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Update(ctx, d, v1.UpdateOptions{})
				if err != nil {
					return controller.NewPermanentError(fmt.Errorf("failed to create service resources: %v", err))
				}
			}

		case schema.FromAPIVersionAndKind("v1", "ConfigMap"):
			logger.Infof("Found ConfigMap: %s", obj.(*corev1.ConfigMap).Name)
			cm := obj.(*corev1.ConfigMap)
			s, _ := c.KubeClientSet.CoreV1().ConfigMaps(p.Namespace).Get(ctx, cm.Name, v1.GetOptions{})
			if s.Name == "" {
				logger.Infof("Create ConfigMap")
				_, err := c.KubeClientSet.CoreV1().ConfigMaps(p.Namespace).Create(ctx, cm, v1.CreateOptions{})
				if err != nil {
					return controller.NewPermanentError(fmt.Errorf("failed to create ConfigMap resources: %v", err))
				}
			} else {
				logger.Infof("Configmap pre-existing, please manually remove configmap %s to apply new changes: %v", cm.Name, s.Name)
			}

		case schema.FromAPIVersionAndKind("v1", "PersistentVolumeClaim"):
			logger.Infof("Found pvc: %s", obj.(*corev1.PersistentVolumeClaim).Name)
			pvc := obj.(*corev1.PersistentVolumeClaim)
			pvcG, _ := c.KubeClientSet.CoreV1().PersistentVolumeClaims(p.Namespace).Get(ctx, pvc.Name, v1.GetOptions{})
			if pvcG.Name == "" {
				logger.Infof("Create pvc: %s", pvc.Name)
				_, err := c.KubeClientSet.CoreV1().PersistentVolumeClaims(p.Namespace).Create(ctx, pvc, v1.CreateOptions{})
				if err != nil {
					return controller.NewPermanentError(fmt.Errorf("failed to create pvc resources: %v", err))
				}
			} else {
				logger.Infof("PVC found, skipping creation")
			}

		default:
			logger.Infof("Type %s was unexpected", obj.GetObjectKind().GroupVersionKind())
			return controller.NewPermanentError(fmt.Errorf("unexpected type"))
		}
	}
	// Status updates
	var bootstrap string
	var testcatalog string
	externalService, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, p.Name+"-external-service", v1.GetOptions{})
	if err != nil {
		logger.Warnf("Problem locating external service: %v", err)
		return controller.NewRequeueAfter(time.Second * 3)
	}
	for _, port := range externalService.Spec.Ports {
		if port.Name == "http" {
			bootstrap = p.Spec.ComponentParms["hostname"] + ":" + strconv.FormatInt(int64(port.NodePort), 10) + "/bootstrap"
			testcatalog = p.Spec.ComponentParms["hostname"] + ":" + strconv.FormatInt(int64(port.NodePort), 10) + "/testcatalog"
		}
	}

	dep, err := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Get(ctx, p.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	if dep.Status.ReadyReplicas == *p.Spec.Replicas {
		p.Status = v2alpha1.ComponentStatus{
			Ready: true,
			StatusParms: map[string]string{
				"bootstrap":   bootstrap,
				"testcatalog": testcatalog,
			},
		}
		c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaApiComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
		return nil
	} else {
		p.Status = v2alpha1.ComponentStatus{
			Ready: false,
		}
		c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaApiComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
		return controller.NewRequeueAfter(time.Second * 3)
	}
}

func mustNewRequirement(key string, op selection.Operator, vals []string) labels.Requirement {
	r, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		panic(fmt.Sprintf("mustNewRequirement(%v, %v, %v) = %v", key, op, vals, err))
	}
	return *r
}
