package local

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
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

// GetURL _
func (s *StorageClaim) GetURL() string {
	return fmt.Sprintf("file://%s", s.file.FullPath())
}

// GetName _
func (s *StorageClaim) GetName() string {
	return s.file.Name()
}

// GetSize _
func (s *StorageClaim) GetSize() int {
	return s.size
}

// WriteFrom _
func (s *StorageClaim) WriteFrom(reader io.Reader) error {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	if s.file == nil {
		return models.ErrNotFound
	}

	fileWriter, err := s.file.WriteContent()

	if err != nil {
		return errors.Wrap(err, "Building file writer")
	}

	_, err = io.Copy(fileWriter, reader)

	fileWriter.Close()

	if err != nil {
		return errors.Wrap(err, "Copying from reader to file writer")
	}

	return nil
}

// ReadTo _
func (s *StorageClaim) ReadTo(writer io.Writer) error {
	if closer, ok := writer.(io.Closer); ok {
		defer closer.Close()
	}

	if s.file == nil {
		return models.ErrNotFound
	}

	fileReader, err := s.file.ReadContent()

	if err != nil {
		return errors.Wrap(err, "Building file reader")
	}

	_, err = io.Copy(writer, fileReader)

	fileReader.Close()

	if err != nil {
		return errors.Wrap(err, "Copying from file reader to writer")
	}

	return nil
}

// // GetReader _
// func (s *StorageClaim) GetReader() (io.ReadCloser, error) {
// 	if s.file == nil {
// 		return nil, ErrStorageClaimMissingFile
// 	}

// 	return s.file.ReadContent()
// }
