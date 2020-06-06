package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailorman/ffchunker/pkg/files"
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
