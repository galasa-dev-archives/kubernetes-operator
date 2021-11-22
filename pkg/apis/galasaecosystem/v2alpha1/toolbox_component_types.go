package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genreconciler:krshapedlogic=false
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GalasaToolboxComponent struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec ToolboxSpec `json:"spec,omitempty"`

	// +optional
	Status ComponentStatus `json:"status,omitempty"`
}

type ToolboxSpec struct {
	// +optional
	Simbank ComponentSpec `json:"simbank"`
	// +optional
	Grafana ComponentSpec `json:"grafana"`
	// +optional
	Prometheus ComponentSpec `json:"prometheus"`
	// +optional
	Jenkins ComponentSpec `json:"Jenkins"`
	// +optional
	Nexus ComponentSpec `json:"Nexus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GalasaToolboxComponentList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GalasaToolboxComponent `json:"items"`
}
