package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Filer _
type Filer interface {
	FullPath() string
	Name() string
	SetChTime(timeObj time.Time) error
	EnsureParentDirExists() error
	Remove() error
	SetDirPath(path Pather)
	SetFileName(fileName string)
	Clone() Filer
	BaseName() string
	Extension() string
	NewWithSuffix(suffix string) Filer
	BuildPath() Pather
	IsExist() bool
	ReadContent() (string, error)
	MarshalYAML() (interface{}, error)
}

// File _
type File struct {
	fileName string
	dirPath  string
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

// BuildPath _
func (f *File) BuildPath() Pather {
	return NewPath(f.DirPath())
}

// IsExist _
func (f *File) IsExist() bool {
	info, err := os.Stat(f.FullPath())
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// SetDirPath _
func (f *File) SetDirPath(path Pather) {
	f.dirPath = path.FullPath()
}

// SetFileName _
func (f *File) SetFileName(fileName string) {
	f.fileName = fileName
}

// Clone _
func (f *File) Clone() Filer {
	newFile := &File{}
	*newFile = *f
	return newFile
}

// BaseName returns file name without extension
func (f *File) BaseName() string {
	return strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
}

// Extension returns file extension from name. Example: ".mp4"
func (f *File) Extension() string {
	return filepath.Ext(f.Name())
}

// NewWithSuffix _
func (f *File) NewWithSuffix(suffix string) Filer {
	newFile := f.Clone()

	nameWithoutExt := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))

	newFile.SetFileName(
		fmt.Sprintf(
			"%s%s%s",
			nameWithoutExt,
			suffix,
			filepath.Ext(f.Name()),
		),
	)

	return newFile
}

// SetChTime _
func (f *File) SetChTime(timeObj time.Time) error {
	return os.Chtimes(f.FullPath(), timeObj, timeObj)
}

// EnsureParentDirExists _
func (f *File) EnsureParentDirExists() error {
	path := NewPath(".")

	return path.Create()
}

// Remove _
func (f *File) Remove() error {
	return os.Remove(f.FullPath())
}

// ReadContent _
func (f *File) ReadContent() (string, error) {
	file, err := os.Open(f.FullPath())

	if err != nil {
		return "", err
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	return string(b), nil
}

// MarshalYAML is YAML Marshaller interface implementation
func (f *File) MarshalYAML() (interface{}, error) {
	return f.FullPath(), nil
}
