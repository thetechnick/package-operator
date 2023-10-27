package packagerepositories

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"package-operator.run/internal/apis/manifests"
	"package-operator.run/internal/packages/internal/packagetypes"
)

func TestPackageIndex(t *testing.T) {
	t.Parallel()
	baseDir := "testdata/packageindex"
	fs := DirFS(baseDir)
	// pi := NewPackageIndex("my-pkg", fs)
	ctx := context.Background()
	pi, err := ReadPackageIndex(ctx, fs, "my-pkg")
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
		ManifestLock: &manifests.PackageManifestLock{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pkg",
			},
		},
		Files: packagetypes.Files{
			"icon.png":         []byte{},
			"xxx/xxx/icon.svg": []byte{},
			"README.md":        []byte{},
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

	// remove entries loaded from disk
	require.NoError(t, pi.Remove(ctx, entry1))
	require.NoError(t, pi.Remove(ctx, entry2))
	require.NoError(t, pi.Tidy(ctx))

	pkgDir := filepath.Join(baseDir, "my-pkg")
	entry1Dir := filepath.Join(pkgDir, "12345")
	entry2Dir := filepath.Join(pkgDir, "67890")
	_, err = os.Stat(entry1Dir)
	require.True(t, os.IsNotExist(err), "my-pkg 12345 folder removed")
	_, err = os.Stat(entry2Dir)
	require.True(t, os.IsNotExist(err), "my-pkg 67890 folder removed")

	// add them back
	require.NoError(t, pi.Add(ctx, entry1, pkg))
	require.NoError(t, pi.Add(ctx, entry2, pkg))
	_, err = os.Stat(filepath.Join(entry1Dir, "icon.png"))
	require.NoError(t, err, "icon.png should be indexed")
	_, err = os.Stat(filepath.Join(entry1Dir, "icon.svg"))
	require.NoError(t, err, "icon.svg should be indexed")
	_, err = os.Stat(filepath.Join(entry1Dir, "readme.md"))
	require.NoError(t, err, "readme.md should be indexed")

	require.NoError(t, pi.Write(ctx))
}
