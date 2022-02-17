/*
 * Copyright contributors to the Galasa Project
 */
package cps

import (
	"context"

	galasacpsreconciler "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasacpscomponent"

	galasaecosystemclient "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/client"
	galasacpsformer "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasacpscomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting CPS controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		cpsinformer := galasacpsformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:            kubeclientset,
			GalasaEcosystemClientSet: clientset,
			GalasaCpsLister:          cpsinformer.Lister(),
		}

		impl := galasacpsreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		cpsinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
