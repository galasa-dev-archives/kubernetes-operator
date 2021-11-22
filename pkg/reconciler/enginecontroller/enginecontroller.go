/*
 * Copyright contributors to the Galasa Project
 */
package enginecontroller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	ecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem"
	galasaecosystem "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	galasaecosystemlisters "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/listers/galasaecosystem/v2alpha1"
	"k8s.io/client-go/kubernetes"
)

type Reconciler struct {
	KubeClientSet kubernetes.Interface

	GalasaEcosystemClientSet     galasaecosystem.Interface
	GalasaEngineControllerLister galasaecosystemlisters.GalasaEngineControllerComponentLister
}

func (c *Reconciler) ReconcileKind(ctx context.Context, p *v2alpha1.GalasaEngineControllerComponent) pkgreconciler.Event {
	logger := logging.FromContext(ctx)
	ec := ecosystem.EngineController(p, c.GalasaEcosystemClientSet)
	objects := ec.GetObjects()
	logger.Infof("Reconciling Engine Controller")

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

		default:
			logger.Infof("Type %s was unexpected", obj.GetObjectKind().GroupVersionKind())
			return controller.NewPermanentError(fmt.Errorf("unexpected type"))
		}
	}
	// Status updates
	dep, err := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Get(ctx, p.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	if dep.Status.ReadyReplicas == *p.Spec.Replicas {
		p.Status = v2alpha1.ComponentStatus{
			Ready: true,
		}
		c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaEngineControllerComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
		return nil
	} else {
		p.Status = v2alpha1.ComponentStatus{
			Ready: false,
		}
		c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaEngineControllerComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
		return controller.NewRequeueAfter(time.Second * 3)
	}
}
