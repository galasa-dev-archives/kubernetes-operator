/*
 * Copyright contributors to the Galasa Project
 */
package toolbox

import (
	"context"

	galasatoolboxreconciler "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasatoolboxcomponent"

	galasaecosystemclient "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/client"
	galasaecosystemformer "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaecosystem"
	galasatoolboxinformer "github.com/galasa-dev/kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasatoolboxcomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting Toolbox controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		informer := galasaecosystemformer.Get(ctx)
		toolboxinformer := galasatoolboxinformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:            kubeclientset,
			GalasaEcosystemClientSet: clientset,
			GalasaEcosystemLister:    informer.Lister(),
			GalasaToolboxLister:      toolboxinformer.Lister(),
		}

		impl := galasatoolboxreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		toolboxinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
