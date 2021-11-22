/*
 * Copyright contributors to the Galasa Project
 */
package enginecontroller

import (
	"context"

	galasaenginecontrollerreconciler "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasaenginecontrollercomponent"

	galasaecosystemclient "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/client"
	galasaenginecontrollerinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaenginecontrollercomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting Engine controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		ecinformer := galasaenginecontrollerinformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:                kubeclientset,
			GalasaEcosystemClientSet:     clientset,
			GalasaEngineControllerLister: ecinformer.Lister(),
		}

		impl := galasaenginecontrollerreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		ecinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
