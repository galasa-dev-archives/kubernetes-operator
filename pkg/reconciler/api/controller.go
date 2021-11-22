/*
 * Copyright contributors to the Galasa Project
 */
package api

import (
	"context"

	galasaapireconciler "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasaapicomponent"

	galasaecosystemclient "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/client"
	galasaapiiformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaapicomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting API controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		apiinformer := galasaapiiformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:            kubeclientset,
			GalasaEcosystemClientSet: clientset,
			GalasaApiLister:          apiinformer.Lister(),
		}

		impl := galasaapireconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		apiinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
