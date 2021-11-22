package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "galasa.dev", Version: "v2alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds Build types to the scheme.
	AddToScheme = schemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&GalasaEcosystem{},
		&GalasaEcosystemList{},
		&GalasaCpsComponent{},
		&GalasaCpsComponentList{},
		&GalasaRasComponent{},
		&GalasaRasComponentList{},
		&GalasaApiComponent{},
		&GalasaApiComponentList{},
		&GalasaEngineControllerComponent{},
		&GalasaEngineControllerComponentList{},
		&GalasaResmonComponent{},
		&GalasaResmonComponentList{},
		&GalasaMetricsComponent{},
		&GalasaMetricsComponentList{},
		&GalasaToolboxComponent{},
		&GalasaToolboxComponentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
