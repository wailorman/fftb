package local

import (
	"io"

	"github.com/wailorman/fftb/pkg/files"
)

// StorageClaim _
type StorageClaim struct {
	identity string
	file     files.Filer
	size     int
}

// GetID _
func (s *StorageClaim) GetID() string {
	return s.identity
}

// GetName _
func (s *StorageClaim) GetName() string {
	return s.file.Name()
}

// GetSize _
func (s *StorageClaim) GetSize() int {
	return s.size
}

// GetWriter _
func (s *StorageClaim) GetWriter() (io.WriteCloser, error) {
	if s.file == nil {
		return nil, ErrStorageClaimMissingFile
	}

	return s.file.WriteContent()
}

// GetReader _
func (s *StorageClaim) GetReader() (io.ReadCloser, error) {
	if s.file == nil {
		return nil, ErrStorageClaimMissingFile
	}

	return s.file.ReadContent()
}
