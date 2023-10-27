package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
type RepositoryPackageIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RepositoryPackageIndexSpec `json:"spec,omitempty"`
}

type RepositoryPackageIndexSpec struct {
	Versions []string                      `json:"versions,omitempty"`
	Index    []RepositoryPackageIndexEntry `json:"index,omitempty"`
}

type RepositoryPackageIndexEntry struct {
	// Versions that reference this digest.
	Versions []string `json:"versions"`
	// Digest of the package image.
	Digest string `json:"digest"`
	// Image of the package. e.g. quay.io/xxx/xxx.
	Image string `json:"image"`
}

func init() { register(&RepositoryPackageIndex{}) }
