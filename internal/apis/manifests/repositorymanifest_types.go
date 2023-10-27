package manifests

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
type RepositoryManifest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RepositoryManifestSpec `json:"spec,omitempty"`
}

type RepositoryManifestSpec struct {
	Index []RepositoryManifestIndexEntry
}

type RepositoryManifestIndexEntry struct {
	Name string
}

// Repository PackageManifest Annotations.
const (
	RepositoryPackageManifestVersionsAnnotation = "repository.package-operator.run/versions"
	RepositoryPackageManifestDigestAnnotation   = "repository.package-operator.run/digest"
	RepositoryPackageManifestImageAnnotation    = "repository.package-operator.run/image"
)

func init() { register(&RepositoryManifest{}) }
