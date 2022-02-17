/*
 * Copyright contributors to the Galasa Project
 */
package toolbox

import (
	"context"
	"fmt"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	ecosystem "github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	galasaecosystem "github.com/galasa-dev/kubernetes-operator/pkg/client/clientset/versioned"
	galasaecosystemlisters "github.com/galasa-dev/kubernetes-operator/pkg/client/listers/galasaecosystem/v2alpha1"
	"k8s.io/client-go/kubernetes"
)

type Reconciler struct {
	KubeClientSet kubernetes.Interface

	GalasaEcosystemClientSet galasaecosystem.Interface
	GalasaEcosystemLister    galasaecosystemlisters.GalasaEcosystemLister
	GalasaToolboxLister      galasaecosystemlisters.GalasaToolboxComponentLister
}

func (c *Reconciler) ReconcileKind(ctx context.Context, p *v2alpha1.GalasaToolboxComponent) pkgreconciler.Event {
	logger := logging.FromContext(ctx)
	toolbox := ecosystem.Toolbox(p, c.GalasaEcosystemClientSet)
	objects := toolbox.GetObjects()
	logger.Infof("Reconciling Toolbox components")

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

		default:
			logger.Infof("Type %s was unexpected", obj.GetObjectKind().GroupVersionKind())
			return controller.NewPermanentError(fmt.Errorf("unexpected type"))
		}
	}
	dep, err := c.KubeClientSet.AppsV1().Deployments(p.Namespace).Get(ctx, p.Name+"-simbank", v1.GetOptions{})
	if err != nil {
		return err
	}
	if dep.Status.ReadyReplicas != 1 {
		p.Status = v2alpha1.ComponentStatus{
			Ready: false,
		}
		c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaToolboxComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
		return controller.NewRequeueAfter(time.Second * 3)
	}

	// Status updates
	var webserviceport string
	var databaseport string
	var zosmfport string
	var telnetport string
	telnetService, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, p.Name+"-simbank-telnet-external", v1.GetOptions{})
	if err != nil {
		logger.Warnf("Problem locating external service: %v", err)
		return controller.NewRequeueAfter(time.Second * 3)
	}
	for _, port := range telnetService.Spec.Ports {
		if port.Name == "simbank-telnet-external" {
			telnetport = strconv.FormatInt(int64(port.NodePort), 10)
		}
	}
	databaseService, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, p.Name+"-simbank-database-external", v1.GetOptions{})
	if err != nil {
		logger.Warnf("Problem locating external service: %v", err)
		return controller.NewRequeueAfter(time.Second * 3)
	}
	for _, port := range databaseService.Spec.Ports {
		if port.Name == "simbank-database-external" {
			databaseport = strconv.FormatInt(int64(port.NodePort), 10)
		}
	}
	webService, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, p.Name+"-simbank-webservice-external", v1.GetOptions{})
	if err != nil {
		logger.Warnf("Problem locating external service: %v", err)
		return controller.NewRequeueAfter(time.Second * 3)
	}
	for _, port := range webService.Spec.Ports {
		if port.Name == "simbank-webservice-external" {
			webserviceport = strconv.FormatInt(int64(port.NodePort), 10)
		}
	}
	zosmfService, err := c.KubeClientSet.CoreV1().Services(p.Namespace).Get(ctx, p.Name+"-simbank-zosmf-external", v1.GetOptions{})
	if err != nil {
		logger.Warnf("Problem locating external service: %v", err)
		return controller.NewRequeueAfter(time.Second * 3)
	}
	for _, port := range zosmfService.Spec.Ports {
		if port.Name == "simbank-zosmf-external" {
			zosmfport = strconv.FormatInt(int64(port.NodePort), 10)
		}
	}
	if zosmfport == "" {
		return controller.NewRequeueAfter(time.Second * 3)
	}
	if webserviceport == "" {
		return controller.NewRequeueAfter(time.Second * 3)
	}
	if databaseport == "" {
		return controller.NewRequeueAfter(time.Second * 3)
	}
	if telnetport == "" {
		return controller.NewRequeueAfter(time.Second * 3)
	}

	p.Status = v2alpha1.ComponentStatus{
		Ready: true,
		StatusParms: map[string]string{
			"simbank-hostname":       p.Spec.SimbankSpec.ComponentParms["hostname"],
			"simbank-telnetport":     telnetport,
			"simbank-databaseport":   databaseport,
			"simbank-webserviceport": webserviceport,
			"simbank-zosmfport":      zosmfport,
		},
	}
	c.GalasaEcosystemClientSet.GalasaV2alpha1().GalasaToolboxComponents(p.Namespace).UpdateStatus(ctx, p, v1.UpdateOptions{})
	return nil
}
