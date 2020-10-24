package files

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Pather _
type Pather interface {
	Files() ([]*Filer, error)
	FullPath() string
	Create() error
}

// Path _
type Path struct {
	path string
}

// NewPath _
func NewPath(relativePath string) *Path {
	fullPath, _ := filepath.Abs(relativePath)

	return &Path{
		path: fullPath,
	}
}

// FullPath _
func (p *Path) FullPath() string {
	return p.path
}

// Create _
func (p *Path) Create() error {
	if _, err := os.Stat(p.FullPath()); os.IsNotExist(err) {
		return os.MkdirAll(p.FullPath(), os.FileMode.Perm(0755))
	}

	return nil
}

// Clone _
func (p *Path) Clone() *Path {
	newPath := &Path{}
	*newPath = *p
	return newPath
}

// Files _
func (p *Path) Files() ([]*Filer, error) {
	files := make([]*File, 0)

	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Walking in path")
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, NewFile(path))

		return nil
	})

	return files, err
}
