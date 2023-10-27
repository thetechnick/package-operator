package manifests

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
type RepositoryPackageIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RepositoryPackageIndexSpec
}

type RepositoryPackageIndexSpec struct {
	Versions []string
	Index    []RepositoryPackageIndexEntry
}

type RepositoryPackageIndexEntry struct {
	// Versions that reference this digest.
	Versions []string
	// Digest of the package image.
	Digest string
	// Image of the package. e.g. quay.io/xxx/xxx.
	Image string
}
