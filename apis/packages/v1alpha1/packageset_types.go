package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackageSetSpec defines the desired state of a PackageSet.
type PackageSetSpec struct {
	// Archived will delete all unpaused objects and remove o
	Archived bool `json:"archived,omitempty"`
	// Paused disables reconcilation of the PackageSet,
	// only Status updates will be propagated.
	Paused bool `json:"paused,omitempty"`
	// Pause reconcilation of specific objects.
	PausedFor []PackagePausedObject `json:"pausedFor,omitempty"`

	// Immutable fields below
	PackageSetTemplateSpec `json:",inline"`
}

type PackageDependencyCRD struct {
	// Name of the Custom Resource Definition to probe for.
	Name string `json:"name"`
	// API Version the CRD must
	Version string `json:"version"`
}

// Specifies that the reconcilation of a specific object should be paused.
type PackagePausedObject struct {
	// Object Kind.
	Kind string `json:"kind"`
	// Object Group.
	Group string `json:"group"`
	// Object Name.
	Name string `json:"name"`
}

// PackageSetStatus defines the observed state of a PackageSet
type PackageSetStatus struct {
	// The most recent generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions is a list of status conditions ths object is in.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// DEPRECATED: This field is not part of any API contract
	// it will go away as soon as kubectl can print conditions!
	// Human readable status - please use .Conditions from code
	Phase PackageSetPhase `json:"phase,omitempty"`
	// List of objects, the controller has paused reconcilation on.
	PausedFor []PackagePausedObject `json:"pausedFor,omitempty"`
}

const (
	PackageSetAvailable = "Available"
	PackageSetPaused    = "Paused"
	PackageSetArchived  = "Archived"
	// Succeeded condition is only set once,
	// after a PackageSet became Available for the first time.
	PackageSetSucceeded = "Succeeded"
)

type PackageSetPhase string

// Well-known PackageSet Phases for printing a Status in kubectl,
// see deprecation notice in PackageSetStatus for details.
const (
	PackageSetPhasePending           PackageSetPhase = "Pending"
	PackageSetPhaseAvailable         PackageSetPhase = "Available"
	PackageSetPhaseNotReady          PackageSetPhase = "NotReady"
	PackageSetPhaseMissingDependency PackageSetPhase = "MissingDependency"
	PackageSetPhaseArchived          PackageSetPhase = "Archived"
)

// PackageSet is the Schema for the PackageSets API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PackageSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PackageSetSpec `json:"spec,omitempty"`
	// +kubebuilder:default={phase:Pending}
	Status PackageSetStatus `json:"status,omitempty"`
}

// PackageSetList contains a list of PackageSets
// +kubebuilder:object:root=true
type PackageSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PackageSet{}, &PackageSetList{})
}
