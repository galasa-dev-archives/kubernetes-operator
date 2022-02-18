/*
 * Copyright contributors to the Galasa Project
 */
package galasaecosystem

import (
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/api"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/cps"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/enginecontroller"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/metrics"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/ras"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/resmon"
	"github.com/galasa-dev/kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1/toolbox"

	galasaecosystem "github.com/galasa-dev/kubernetes-operator/pkg/client/clientset/versioned"
)

func Cps(cpsCrd *v2alpha1.GalasaCpsComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return cps.New(cpsCrd, k)
}

func Ras(rasCrd *v2alpha1.GalasaRasComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return ras.New(rasCrd, k)
}

func Api(apiCrd *v2alpha1.GalasaApiComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return api.New(apiCrd, k)
}

func Metrics(metricsCrd *v2alpha1.GalasaMetricsComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return metrics.New(metricsCrd, k)
}

func Resmon(resmonCrd *v2alpha1.GalasaResmonComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return resmon.New(resmonCrd, k)
}

func EngineController(engineControllerCrd *v2alpha1.GalasaEngineControllerComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return enginecontroller.New(engineControllerCrd, k)
}

func Toolbox(toolboxCrd *v2alpha1.GalasaToolboxComponent, k galasaecosystem.Interface) v2alpha1.ComponentInterface {
	return toolbox.New(toolboxCrd, k)
}
