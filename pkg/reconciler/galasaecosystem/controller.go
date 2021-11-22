package galasaecosystem

import (
	"context"

	galasaecosystemreconciler "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/reconciler/galasaecosystem/v2alpha1/galasaecosystem"

	galasaecosystemclient "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/client"
	galasaapiinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaapicomponent"
	galasacpsinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasacpscomponent"
	galasaecosystemformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaecosystem"
	galasaenginecontrollerinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaenginecontrollercomponent"
	galasametricsinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasametricscomponent"
	galasarasinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasarascomponent"
	galasaresmoninformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaresmoncomponent"
	galasatoolboxinformer "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasatoolboxcomponent"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

func NewController(namespace string) func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		logger := logging.FromContext(ctx)
		logger.Info("Starting Ecosystem operator...")

		kubeclientset := kubeclient.Get(ctx)
		clientset := galasaecosystemclient.Get(ctx)
		ecosysteminformer := galasaecosystemformer.Get(ctx)
		cpsinformer := galasacpsinformer.Get(ctx)
		rasinformer := galasarasinformer.Get(ctx)
		apiinformer := galasaapiinformer.Get(ctx)
		resmoninformer := galasaresmoninformer.Get(ctx)
		metricsinformer := galasametricsinformer.Get(ctx)
		enginecontrollerinformer := galasaenginecontrollerinformer.Get(ctx)
		toolboxinformer := galasatoolboxinformer.Get(ctx)

		c := &Reconciler{
			KubeClientSet:                kubeclientset,
			GalasaEcosystemClientSet:     clientset,
			GalasaEcosystemLister:        ecosysteminformer.Lister(),
			GalasaCPSLister:              cpsinformer.Lister(),
			GalasaRASLister:              rasinformer.Lister(),
			GalasaAPILister:              apiinformer.Lister(),
			GalasaResmonLister:           resmoninformer.Lister(),
			GalasaEngineControllerLister: enginecontrollerinformer.Lister(),
			GalasaMetricsLister:          metricsinformer.Lister(),
			GalasaToolboxLister:          toolboxinformer.Lister(),
		}

		impl := galasaecosystemreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				ConfigStore: nil,
			}
		})

		ecosysteminformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		cpsinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		rasinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		apiinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		enginecontrollerinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		resmoninformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		toolboxinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
		metricsinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

		return impl
	}
}
