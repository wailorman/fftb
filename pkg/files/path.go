package files

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Pather _
type Pather interface {
	Files() ([]Filer, error)
	FullPath() string
	Create() error
	BuildSubpath(path string) Pather
	BuildFile(filePath string) Filer
	Destroy() error
}

// Path _
type Path struct {
	path string
}

// NewPath _
func NewPath(relativePath string) Pather {
	fullPath, _ := filepath.Abs(relativePath)

	return &Path{
		path: fullPath,
	}
}

// NewTempPath _
func NewTempPath(name string) Pather {
	return NewPath(filepath.Clean(os.TempDir() + "/" + name))
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

// BuildSubpath returns instance of new directory in current directory (not written to disk)
func (p *Path) BuildSubpath(path string) Pather {
	return NewPath(p.FullPath() + "/" + path)
}

// BuildFile returns instance of new file inside current directory (not written to disk)
func (p *Path) BuildFile(filePath string) Filer {
	return NewFile(p.FullPath() + "/" + filePath)
}

// Destroy _
func (p *Path) Destroy() error {
	return os.RemoveAll(p.FullPath())
}

// Clone _
func (p *Path) Clone() *Path {
	newPath := &Path{}
	*newPath = *p
	return newPath
}

// Files _
func (p *Path) Files() ([]Filer, error) {
	files := make([]Filer, 0)

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
