package packagerepositories

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"package-operator.run/internal/apis/manifests"
	"package-operator.run/internal/packages/internal/packagetypes"
)

const (
	packageIndexFilename        = "index.yaml"
	packageManifestFilename     = "manifest.yaml"
	packageManifestLockFilename = "manifest.lock.yaml"
)

type PackageIndex struct {
	fsys     FS
	manifest *manifests.RepositoryPackageIndex

	orderedVersions semver.Collection
	versions        map[string]struct{}
	versionToDigest map[string]string
	digestToMeta    map[string]manifests.RepositoryPackageIndexEntry
}

func newPackageIndex(fsys FS) *PackageIndex {
	return &PackageIndex{
		fsys: fsys,

		versions:        map[string]struct{}{},
		versionToDigest: map[string]string{},
		digestToMeta:    map[string]manifests.RepositoryPackageIndexEntry{},
	}
}

func NewPackageIndex(name string, fsys FS) *PackageIndex {
	pi := newPackageIndex(fsys)
	pi.manifest = &manifests.RepositoryPackageIndex{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return pi
}

func ReadPackageIndex(ctx context.Context, fsys FS, pkgName string) (index *PackageIndex, err error) {
	fbytes, err := fs.ReadFile(fsys, filepath.Join(pkgName, packageIndexFilename))
	if err != nil {
		return
	}

	index = newPackageIndex(fsys)
	index.manifest = &manifests.RepositoryPackageIndex{}
	if err = yaml.Unmarshal(fbytes, index.manifest); err != nil {
		return
	}
	for _, entry := range index.manifest.Spec.Index {
		if err := index.add(ctx, entry); err != nil {
			return index, err
		}
	}
	return
}

func (pi *PackageIndex) IsEmpty() bool {
	return len(pi.digestToMeta) == 0
}

func (pi *PackageIndex) Remove(
	ctx context.Context,
	entry manifests.RepositoryPackageIndexEntry,
) error {
	var (
		orderedVersions  semver.Collection
		versionsToRemove = map[string]struct{}{}
	)
	for _, v := range entry.Versions {
		delete(pi.versions, v)
		delete(pi.versionToDigest, v)
		versionsToRemove[v] = struct{}{}
	}
	for _, sv := range pi.orderedVersions {
		if _, remove := versionsToRemove["v"+sv.String()]; remove {
			continue
		}
		orderedVersions = append(orderedVersions, sv)
	}
	pi.orderedVersions = orderedVersions
	delete(pi.digestToMeta, entry.Digest)
	return nil
}

func (pi *PackageIndex) Add(
	ctx context.Context,
	entry manifests.RepositoryPackageIndexEntry,
	pkg *packagetypes.Package,
) error {
	if pkg.Manifest.Name != pi.manifest.Name {
		return fmt.Errorf("package index for package %s, got: %s", pi.manifest.Name, pkg.Manifest.Name)
	}

	manifestYaml, err := yaml.Marshal(pkg.Manifest)
	if err != nil {
		return err
	}
	path := filepath.Join(pkg.Manifest.Name, entry.Digest, packageManifestFilename)
	if err := pi.fsys.WriteFile(path, manifestYaml); err != nil {
		return err
	}
	if pkg.ManifestLock != nil {
		manifestLockYaml, err := yaml.Marshal(pkg.Manifest)
		if err != nil {
			return err
		}
		if err := pi.fsys.WriteFile(path, manifestLockYaml); err != nil {
			return err
		}
	}
	if err := includeFiles(pi.fsys, entry, pkg); err != nil {
		return err
	}
	return pi.add(ctx, entry)
}

func (pi *PackageIndex) add(
	ctx context.Context,
	entry manifests.RepositoryPackageIndexEntry,
) error {
	var entryOrderedVersions semver.Collection
	for _, v := range entry.Versions {
		if _, ok := pi.versions[v]; ok {
			return fmt.Errorf("version %s already indexed", v)
		}

		sv, err := semver.StrictNewVersion(strings.TrimPrefix(v, "v"))
		if err != nil {
			return err
		}
		pi.versions[v] = struct{}{}
		entryOrderedVersions = append(entryOrderedVersions, sv)
		pi.orderedVersions = append(pi.orderedVersions, sv)
		pi.versionToDigest[v] = entry.Digest
	}
	sort.Sort(pi.orderedVersions)
	slices.Reverse(pi.orderedVersions)
	sort.Sort(entryOrderedVersions)
	slices.Reverse(entryOrderedVersions)

	entry.Versions = nil
	for _, sv := range entryOrderedVersions {
		entry.Versions = append(entry.Versions, "v"+sv.String())
	}
	pi.digestToMeta[entry.Digest] = entry
	return nil
}

// Write the updated PackageIndex to the given FS.
func (pi *PackageIndex) Write(ctx context.Context) error {
	pi.manifest.Spec.Index = nil
	for _, entry := range pi.digestToMeta {
		pi.manifest.Spec.Index = append(pi.manifest.Spec.Index, entry)
	}
	pi.manifest.Spec.Versions = nil
	for _, v := range pi.orderedVersions {
		pi.manifest.Spec.Versions = append(pi.manifest.Spec.Versions, "v"+v.String())
	}

	pi.manifest.CreationTimestamp = metav1.Now()
	indexBytes, err := yaml.Marshal(pi.manifest)
	if err != nil {
		return err
	}
	return pi.fsys.WriteFile(filepath.Join(pi.manifest.Name, packageIndexFilename), indexBytes)
}

// Tidy up all unknown data in the filesystem.
func (pi *PackageIndex) Tidy(ctx context.Context) error {
	dirs, err := fs.ReadDir(pi.fsys, pi.manifest.Name)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		if _, ok := pi.digestToMeta[dir.Name()]; ok {
			continue
		}
		path := filepath.Join(pi.manifest.Name, dir.Name())
		if err := pi.fsys.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

const (
	iconFilename        = "icon"
	readmeFileName      = "readme.md"
	readmeUpperFileName = "README.md"
)

func includeFiles(fsys FS, entry manifests.RepositoryPackageIndexEntry, pkg *packagetypes.Package) error {
	// readme/README.md
	if rm, ok := pkg.Files[readmeFileName]; ok {
		path := filepath.Join(pkg.Manifest.Name, entry.Digest, readmeFileName)
		if err := fsys.WriteFile(path, rm); err != nil {
			return err
		}
	}
	if rm, ok := pkg.Files[readmeUpperFileName]; ok {
		path := filepath.Join(pkg.Manifest.Name, entry.Digest, readmeFileName)
		if err := fsys.WriteFile(path, rm); err != nil {
			return err
		}
	}
	for path, content := range pkg.Files {
		file := filepath.Base(path)
		if strings.HasPrefix(file, iconFilename+".") {
			path := filepath.Join(pkg.Manifest.Name, entry.Digest, file)
			if err := fsys.WriteFile(path, content); err != nil {
				return err
			}
		}
	}
	return nil
}
