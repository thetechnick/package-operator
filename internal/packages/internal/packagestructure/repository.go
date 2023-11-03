package packagestructure

import (
	"context"

	manifestsv1alpha1 "package-operator.run/apis/manifests/v1alpha1"
	"package-operator.run/internal/apis/manifests"
)

func RepositoryPackageIndexFromFile(
	ctx context.Context, path string, manifestBytes []byte,
) (*manifests.RepositoryPackageIndex, error) {
	return manifestFromFile[manifests.RepositoryPackageIndex](ctx, scheme, path, manifestBytes)
}

// Converts the internal version of an RepositoryPackageIndex into it's v1alpha1 representation.
func ToV1Alpha1RepositoryPackageIndex(in *manifests.RepositoryPackageIndex) (*manifestsv1alpha1.RepositoryPackageIndex, error) {
	out := &manifestsv1alpha1.RepositoryPackageIndex{}
	if err := scheme.Convert(in, out, nil); err != nil {
		return nil, err
	}
	out.SetGroupVersionKind(manifestsv1alpha1.GroupVersion.WithKind("RepositoryPackageIndex"))
	return out, nil
}

func RepositoryManifestFromFile(
	ctx context.Context, path string, manifestBytes []byte,
) (*manifests.RepositoryManifest, error) {
	return manifestFromFile[manifests.RepositoryManifest](ctx, scheme, path, manifestBytes)
}

// Converts the internal version of an RepositoryManifest into it's v1alpha1 representation.
func ToV1Alpha1RepositoryManifest(in *manifests.RepositoryManifest) (*manifestsv1alpha1.RepositoryManifest, error) {
	out := &manifestsv1alpha1.RepositoryManifest{}
	if err := scheme.Convert(in, out, nil); err != nil {
		return nil, err
	}
	out.SetGroupVersionKind(manifestsv1alpha1.GroupVersion.WithKind("RepositoryManifest"))
	return out, nil
}
