package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
type RepositoryManifest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RepositoryManifestSpec `json:"spec,omitempty"`
}

type RepositoryManifestSpec struct {
	Index []RepositoryManifestIndexEntry `json:"index,omitempty"`
}

type RepositoryManifestIndexEntry struct {
	Name string `json:"name"`
}

func init() { register(&RepositoryManifest{}) }
