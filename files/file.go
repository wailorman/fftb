package files

import (
	"os"
	"path/filepath"
	"time"
)

// Filer _
type Filer interface {
	FullPath() string
	Name() string
	SetChTime(timeObj time.Time) error
	EnsureParentDirExists() error
	Remove() error
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

// EnsureParentDirExists _
func (f *File) EnsureParentDirExists() error {
	path := NewPathBuilder(f.DirPath()).NewPath(".")

	return path.Create()
}
