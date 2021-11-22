package v2alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genreconciler:krshapedlogic=false
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GalasaEcosystem struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GalasaEcosystemSpec `json:"spec,omitempty"`

	// +optional
	Status GalasaEcosystemStatus `json:"status,omitempty"`
}

type GalasaEcosystemSpec struct {
	Hostname      string `json:"hostname"`
	GalasaVersion string `json:"galasaVersion"`
	// +optional
	BusyboxImage   string                   `json:"busyboxImage"`
	ComponentsSpec map[string]ComponentSpec `json:"componentsSpec"`
}

type ComponentSpec struct {
	Image    string `json:"image"`
	Replicas *int32 `json:"replicas"`
	// +optional
	ImagePullPolicy string `json:"imagePullPolicy"`
	// +optional
	Storage string `json:"storage"`
	// +optional
	StorageClassName string `json:"storageClassName"`
	// +optional
	NodeSelector map[string]string `json:"nodeSelector"`
	// +optional
	ComponentParms map[string]string `json:"componentParms"`
}

type ComponentStatus struct {
	Ready       bool              `json:"ready"`
	StatusParms map[string]string `json:"statusParms"`
}

type ComponentInterface interface {
	IsReady(context.Context) bool

	HasChanged(spec ComponentSpec) bool

	GetObjects() []runtime.Object
}

type GalasaEcosystemStatus struct {
	Ready        bool   `json:"ready"`
	BootstrapURL string `json:"bootstrapURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GalasaEcosystemList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GalasaEcosystem `json:"items"`
}
