/*
 * Copyright contributors to the Galasa Project
 */
package ras

import (
	"context"

	galasarasreconciler "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasarascomponent"

	galasaecosystemclient "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/client"
	galasarasormer "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasarascomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting RAS controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		rasinformer := galasarasormer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:            kubeclientset,
			GalasaEcosystemClientSet: clientset,
			GalasaRasLister:          rasinformer.Lister(),
		}

		impl := galasarasreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		rasinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
