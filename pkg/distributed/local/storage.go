package local

import (
	"io"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
)

// ErrStorageClaimMissingFile _
var ErrStorageClaimMissingFile = errors.New("Storage claim missing file")

// ErrStorageClaimAlreadyAllocated _
var ErrStorageClaimAlreadyAllocated = errors.New("Storage claim already allocated")

// StorageControl _
type StorageControl struct {
	storagePath files.Pather
}

// NewStorageControl _
func NewStorageControl(path files.Pather) *StorageControl {
	return &StorageControl{
		storagePath: path,
	}
}

// AllocateStorageClaim _
func (sc *StorageControl) AllocateStorageClaim(identity string) (models.IStorageClaim, error) {
	file := sc.storagePath.BuildFile(identity)

	err := file.EnsureParentDirExists()

	if err != nil {
		return nil, errors.Wrap(err, "Creating directory for storage claim")
	}

	if file.IsExist() {
		return nil, ErrStorageClaimAlreadyAllocated
	}

	err = file.Create()

	if err != nil {
		return nil, errors.Wrap(err, "Creating file for storage claim")
	}

	claim := &StorageClaim{
		identity: identity,
		file:     file,
	}

	return claim, nil
}

// BuildStorageClaim _
func (sc *StorageControl) BuildStorageClaim(identity string) (models.IStorageClaim, error) {
	claimFile := sc.storagePath.BuildFile(identity)

	if claimFile.IsExist() == false {
		return nil, ErrStorageClaimMissingFile
	}

	return &StorageClaim{
		identity: identity,
		file:     claimFile,
	}, nil
}

// PurgeStorageClaim _
func (sc *StorageControl) PurgeStorageClaim(claim models.IStorageClaim) error {
	localClaim, ok := claim.(*StorageClaim)

	if !ok {
		return models.ErrUnknownStorageClaimType
	}

	if localClaim.file == nil {
		return ErrStorageClaimMissingFile
	}

	err := localClaim.file.Remove()

	if err != nil {
		return errors.Wrap(err, "Removing file")
	}

	return nil
}

// StorageClaim _
type StorageClaim struct {
	identity string
	file     files.Filer
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
func (s *StorageClaim) GetSize() (int, error) {
	if s.file == nil {
		return 0, ErrStorageClaimMissingFile
	}

	return s.file.Size()
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
