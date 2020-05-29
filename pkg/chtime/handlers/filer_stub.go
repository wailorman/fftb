package handlers

import (
	"os"
	"path/filepath"
	"time"
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
