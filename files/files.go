package files

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// Filer _
type Filer interface {
	FullPath() string
	Name() string
	SetChTime(timeObj time.Time) error
}

// Pather _
type Pather interface {
	Files() ([]*File, error)
}

// PathBuilder _
type PathBuilder struct {
	pwd string
}

// NewPathBuilder _
func NewPathBuilder(pwd string) *PathBuilder {
	return &PathBuilder{
		pwd: pwd,
	}
}

// NewFile _
func (pb *PathBuilder) NewFile(relativePath string) *File {
	fullPath := filepath.Join(pb.pwd, relativePath)

	dirPath, fileName := filepath.Split(fullPath)

	return &File{
		pathBuilder: pb,
		fileName:    fileName,
		dirPath:     dirPath,
	}
}

// NewPath _
func (pb *PathBuilder) NewPath(relativePath string) *Path {
	fullPath := filepath.Join(pb.pwd, relativePath)

	return &Path{
		pathBuilder: pb,
		path:        fullPath,
	}
}

// File _
type File struct {
	pathBuilder *PathBuilder
	fileName    string
	dirPath     string
}

// FullPath _
func (f *File) FullPath() string {
	return filepath.Join(f.dirPath, f.fileName)
}

// Name _
func (f *File) Name() string {
	return f.fileName
}

// DirPath _
func (f *File) DirPath() string {
	return f.dirPath
}

// SetChTime _
func (f *File) SetChTime(timeObj time.Time) error {
	return os.Chtimes(f.FullPath(), timeObj, timeObj)
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
