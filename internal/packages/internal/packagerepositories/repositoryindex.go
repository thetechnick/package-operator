package packagerepositories

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"package-operator.run/internal/apis/manifests"
	"package-operator.run/internal/packages/internal/packagestructure"
	"package-operator.run/internal/packages/internal/packagetypes"
)

const (
	repositoryManifestFilename = "repository.yaml"
)

type RepositoryIndex struct {
	fsys     FS
	manifest *manifests.RepositoryManifest

	packageIndexes map[string]*PackageIndex
}

func newRepositoryIndex(fsys FS) *RepositoryIndex {
	return &RepositoryIndex{
		fsys: fsys,

		packageIndexes: map[string]*PackageIndex{},
	}
}

func NewRepositoryIndex(name string, fsys FS) *RepositoryIndex {
	pi := newRepositoryIndex(fsys)
	pi.manifest = &manifests.RepositoryManifest{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return pi
}

func ReadRepositoryIndex(ctx context.Context, fsys FS) (index *RepositoryIndex, err error) {
	fbytes, err := fs.ReadFile(fsys, repositoryManifestFilename)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}

	index = newRepositoryIndex(fsys)
	index.manifest, err = packagestructure.RepositoryManifestFromFile(ctx, repositoryManifestFilename, fbytes)
	if err != nil {
		return nil, err
	}

	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		pi, err := ReadPackageIndex(ctx, fsys, dir.Name())
		if err != nil {
			return nil, err
		}
		index.packageIndexes[dir.Name()] = pi
	}
	return
}

func (ri *RepositoryIndex) Remove(
	ctx context.Context,
	entry manifests.RepositoryPackageIndexEntry,
	pkg *packagetypes.Package,
) error {
	pi, exists := ri.packageIndexes[pkg.Manifest.Name]
	if !exists {
		return nil
	}
	if err := pi.Remove(ctx, entry); err != nil {
		return err
	}
	if pi.IsEmpty() {
		return ri.fsys.RemoveAll(pkg.Manifest.Name)
	}
	return pi.Write(ctx)
}

func (ri *RepositoryIndex) Add(
	ctx context.Context,
	entry manifests.RepositoryPackageIndexEntry,
	pkg *packagetypes.Package,
) error {
	// TODO validation

	pi, exists := ri.packageIndexes[pkg.Manifest.Name]
	if !exists {
		pi = NewPackageIndex(pkg.Manifest.Name, ri.fsys)
		ri.packageIndexes[pkg.Manifest.Name] = pi
	}
	if err := pi.Add(ctx, entry, pkg); err != nil {
		return err
	}
	return pi.Write(ctx)
}

// Write the updated PackageIndex to the given FS.
func (ri *RepositoryIndex) Write(ctx context.Context) error {
	ri.manifest.CreationTimestamp = metav1.Now()

	ri.manifest.Spec.Index = nil
	var pkgNames []string
	for pkgName := range ri.packageIndexes {
		pkgNames = append(pkgNames, pkgName)
	}
	slices.Sort(pkgNames)
	for _, pkgName := range pkgNames {
		ri.manifest.Spec.Index = append(ri.manifest.Spec.Index, manifests.RepositoryManifestIndexEntry{
			Name: pkgName,
		})
	}

	v1alpha1Manifest, err := packagestructure.ToV1Alpha1RepositoryManifest(ri.manifest)
	if err != nil {
		return err
	}
	indexBytes, err := yaml.Marshal(v1alpha1Manifest)
	if err != nil {
		return err
	}
	return ri.fsys.WriteFile(repositoryManifestFilename, indexBytes)
}
