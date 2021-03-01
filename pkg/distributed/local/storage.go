package local

import (
	"context"
	"io"

	"github.com/machinebox/progress"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/ctxio"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
	// "github.com/machinebox/progress"
	// "github.com/machinebox/progress"
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

// UploadFileToStorageClaim _
func UploadFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim) (chan models.Progresser, chan error) {
	p := make(chan models.Progresser)
	failures := make(chan error)

	go func() {
		defer close(p)
		defer close(failures)

		storageWriter, err := sc.GetWriter()

		if err != nil {
			failures <- errors.Wrap(err, "Building storage claim writer")
			return
		}

		fileReader, err := file.ReadContent()

		if err != nil {
			failures <- errors.Wrap(err, "Building file reader")
			return
		}

		ctxReader := ctxio.NewReader(ctx, fileReader)

		progressReader := progress.NewReader(ctxReader)

		_, err = io.Copy(storageWriter, progressReader)

		storageWriter.Close()
		fileReader.Close()

		if err != nil {
			failures <- errors.Wrap(err, "Failed to upload")
			return
		}
	}()

	return p, failures
}

// DownloadFileFromStorageClaim _
func DownloadFileFromStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim) (chan models.Progresser, chan error) {
	p := make(chan models.Progresser)
	failures := make(chan error)

	go func() {
		defer close(p)
		defer close(failures)

		fileWriter, err := file.WriteContent()

		if err != nil {
			failures <- errors.Wrap(err, "Building file writer")
			return
		}

		storageReader, err := sc.GetReader()

		if err != nil {
			failures <- errors.Wrap(err, "Building storage claim reader")
			return
		}

		// ctxReader := ctxio.NewReader(ctx, storageReader)

		// progressReader := progress.NewReader(ctxReader)

		// scSize, err := sc.GetSize()

		// if err != nil {
		// 	failures <- errors.Wrap(err, "Getting storage claim size")
		// 	return
		// }

		// go func() {
		// 	progressChan := progress.NewTicker(ctx, progressReader, int64(scSize), 1*time.Second)
		// 	for pM := range progressChan {
		// 		p <- dlog.MakeIOProgress("downloading", pM.Percent())
		// 	}
		// }()

		_, err = io.Copy(fileWriter, storageReader)
		// _, err = io.Copy(fileWriter, progressReader)

		fileWriter.Close()
		storageReader.Close()

		if err != nil {
			failures <- errors.Wrap(err, "Failed to download")
			return
		}
	}()

	return p, failures
}
