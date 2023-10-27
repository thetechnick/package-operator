package packagerepositories

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FS interface {
	fs.FS
	WriteFile(name string, content []byte) error
	RemoveAll(name string) error
}

func DirFS(path string) FS {
	return &dirFS{
		FS:   os.DirFS(path),
		path: path,
	}
}

type dirFS struct {
	fs.FS
	path string
}

func (fs *dirFS) WriteFile(name string, content []byte) error {
	path := filepath.Join(fs.path, name)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return os.WriteFile(filepath.Join(fs.path, name), content, os.ModePerm)
}

func (fs *dirFS) RemoveAll(name string) error {
	return os.RemoveAll(filepath.Join(fs.path, name))
}
