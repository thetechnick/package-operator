package packagerepositories

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"package-operator.run/internal/apis/manifests"
	"package-operator.run/internal/packages/internal/packagetypes"
)

func TestRepositoryIndex(t *testing.T) {
	t.Parallel()
	fs := DirFS("testdata/repo")
	ctx := context.Background()

	ri, err := ReadRepositoryIndex(ctx, fs)
	require.NoError(t, err)

	pkg := &packagetypes.Package{
		Manifest: &manifests.PackageManifest{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pkg",
			},
			Repository: manifests.PackageManifestRepository{
				DisplayName: "banana",
			},
		},
	}
	entry1 := manifests.RepositoryPackageIndexEntry{
		Digest:   "12345",
		Image:    "quay.io/example/example",
		Versions: []string{"v1.0.0", "v1.0.0-rc.1"},
	}
	entry2 := manifests.RepositoryPackageIndexEntry{
		Digest:   "67890",
		Image:    "quay.io/example/example",
		Versions: []string{"v1.1.0", "v1.1.0-rc.1"},
	}
	require.NoError(t, ri.Remove(ctx, entry1, pkg))
	_, err = os.Stat("testdata/repo/my-pkg")
	require.NoError(t, err, "my-pkg folder still exists")

	require.NoError(t, ri.Remove(ctx, entry2, pkg))
	_, err = os.Stat("testdata/repo/my-pkg")
	require.True(t, os.IsNotExist(err), "my-pkg folder removed")

	require.NoError(t, ri.Add(ctx, entry1, pkg))
	require.NoError(t, ri.Add(ctx, entry2, pkg))
	require.NoError(t, ri.Write(ctx))
}
