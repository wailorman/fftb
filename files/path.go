package files

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Pather _
type Pather interface {
	Files() ([]*File, error)
	FullPath() string
	Create() error
}

// Path _
type Path struct {
	pathBuilder *PathBuilder
	path        string
}

// Files _
func (p *Path) Files() ([]*File, error) {
	files := make([]*File, 0)

	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Walking in path")
		}

		if info.IsDir() {
			return nil
		}

		dir, name := filepath.Split(path)

		files = append(files, NewPathBuilder(dir).NewFile(name))

		return nil
	})

	return files, err
}
