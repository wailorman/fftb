package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

// filerStub _
type filerStub struct {
	fileName string
	dirPath  string
}

func newFilerStub(dirPath, fileName string) *filerStub {
	return &filerStub{
		fileName: fileName,
		dirPath:  dirPath,
	}
}

// FullPath _
func (f *filerStub) FullPath() string {
	return filepath.Join(f.dirPath, f.fileName)
}

// Name _
func (f *filerStub) Name() string {
	return f.fileName
}

// DirPath _
func (f *filerStub) DirPath() string {
	return f.dirPath
}

// SetChTime _
func (f *filerStub) SetChTime(timeObj time.Time) error {
	return nil
}

// EnsureParentDirExists
func (f *filerStub) EnsureParentDirExists() error {
	return nil
}

// Remove
func (f *filerStub) Remove() error {
	return os.Remove(f.FullPath())
}

// SetFileName _
func (f *filerStub) SetFileName(fileName string) {
	f.fileName = fileName
}

// SetDirPath _`
func (f *filerStub) SetDirPath(path files.Pather) {
	f.dirPath = path.FullPath()
}

// BuildPath _
func (f *filerStub) BuildPath() files.Pather {
	return files.NewPath(f.DirPath())
}

// IsExist _
func (f *filerStub) IsExist() bool {
	return false
}

/// Clone _
func (f *filerStub) Clone() files.Filer {
	newFile := &filerStub{}
	*newFile = *f
	return newFile
}

// Move _
func (f *filerStub) Move(newFullPath string) error {
	return nil
}

// Rename _
func (f *filerStub) Rename(newName string) error {
	return nil
}

// BaseName returns file name without extension
func (f *filerStub) BaseName() string {
	return strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
}

// Extension returns file extension from name. Example: ".mp4"
func (f *filerStub) Extension() string {
	return filepath.Ext(f.Name())
}

// NewWithSuffix _
func (f *filerStub) NewWithSuffix(suffix string) files.Filer {
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

// ReadAllContent _
func (f *filerStub) ReadAllContent() (string, error) {
	file, err := os.Open(f.FullPath())

	if err != nil {
		return "", err
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	return string(b), nil
}

// ReadContent _
func (f *filerStub) ReadContent() (files.FileReader, error) {
	return os.Open(f.FullPath())
}

// Create _
func (f *filerStub) Create() error {
	return nil
}

// Size _
func (f *filerStub) Size() (int, error) {
	info, err := os.Stat(f.FullPath())

	if err != nil {
		return 0, errors.Wrap(err, "Getting file size")
	}

	return int(info.Size()), nil
}

// WriteContent _
func (f *filerStub) WriteContent() (files.FileWriter, error) {
	return nil, nil
}

// MarshalYAML is YAML Marshaller interface implementation
func (f *filerStub) MarshalYAML() (interface{}, error) {
	return f.FullPath(), nil
}
