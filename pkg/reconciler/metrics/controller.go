/*
 * Copyright contributors to the Galasa Project
 */
package metrics

import (
	"context"

	galasametricsreconciler "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasametricscomponent"

	galasaecosystemclient "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/client"
	galasametricsinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasametricscomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting Metrics controller...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		metricsinformer := galasametricsinformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:            kubeclientset,
			GalasaEcosystemClientSet: clientset,
			GalasaMetricsLister:      metricsinformer.Lister(),
		}

		impl := galasametricsreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		metricsinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
