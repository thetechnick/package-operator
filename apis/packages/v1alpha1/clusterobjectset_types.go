package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterObjectSetSpec defines the desired state of a ClusterObjectSet.
type ClusterObjectSetSpec struct {
	// Specifies the lifecycle state of the ObjectSet.
	// +kubebuilder:default="Active"
	// +kubebuilder:validation:Enum=Active;Paused;Archived
	LifecycleState ObjectSetLifecycleState `json:"lifecycleState,omitempty"`
	// Pause reconcilation of specific objects, while still reporting status.
	PausedFor []ObjectSetPausedObject `json:"pausedFor,omitempty"`
	// Immutable fields below
	ObjectSetTemplateSpec `json:",inline"`
}

// ClusterObjectSetStatus defines the observed state of a ClusterObjectSet
type ClusterObjectSetStatus struct {
	// Conditions is a list of status conditions ths object is in.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// DEPRECATED: This field is not part of any API contract
	// it will go away as soon as kubectl can print conditions!
	// Human readable status - please use .Conditions from code
	Phase ObjectSetStatusPhase `json:"phase,omitempty"`
	// List of objects, the controller has paused reconcilation on.
	PausedFor []ObjectSetPausedObject `json:"pausedFor,omitempty"`
}

// ClusterObjectSet reconcile a collection of objects across ordered phases and aggregate their status.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClusterObjectSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterObjectSetSpec `json:"spec,omitempty"`
	// +kubebuilder:default={phase:Pending}
	Status ClusterObjectSetStatus `json:"status,omitempty"`
}

// ClusterObjectSetList contains a list of ClusterObjectSets
// +kubebuilder:object:root=true
type ClusterObjectSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterObjectSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterObjectSet{}, &ClusterObjectSetList{})
}
