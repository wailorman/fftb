package local

import (
	"context"
	"io"

	"github.com/machinebox/progress"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/ctxio"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
)

// StorageClient _
type StorageClient struct {
	localCopiesPath string
}

// NewStorageClient _
func NewStorageClient(localCopiesPath string) *StorageClient {
	return &StorageClient{
		localCopiesPath: localCopiesPath,
	}
}

// RemoveLocalCopy _
func (client *StorageClient) RemoveLocalCopy(ctx context.Context, sc models.IStorageClaim) error {
	// if isLocalStorageClaim(sc) {
	// 	return nil
	// }

	copyFile := buildLocalCopyFile(client.localCopiesPath, sc)

	if !copyFile.IsExist() {
		return nil
	}

	return copyFile.Remove()
}

// MakeLocalCopy _
func (client *StorageClient) MakeLocalCopy(ctx context.Context, sc models.IStorageClaim, p chan models.Progresser) (files.Filer, error) {
	copyFile := buildLocalCopyFile(client.localCopiesPath, sc)

	if client.isLocalCopyMatches(sc) {
		return copyFile, nil
	}

	// if isLocalStorageClaim(sc) {
	// 	return nil
	// }

	if err := copyFile.Create(); err != nil {
		return nil, errors.Wrap(err, "Creating & truncating file for storage claim local copy")
	}

	err := client.DownloadFileFromStorageClaim(ctx, copyFile, sc, p)

	if err != nil {
		return nil, errors.Wrap(err, "Downloading storage claim to local copy")
	}

	return copyFile, nil
}

func (client *StorageClient) isLocalCopyMatches(sc models.IStorageClaim) bool {
	copyFile := buildLocalCopyFile(client.localCopiesPath, sc)

	if !copyFile.IsExist() {
		return false
	}

	copyFileSize, err := copyFile.Size()

	if err != nil {
		return false
	}

	return copyFileSize == sc.GetSize()
}

func buildLocalCopyFile(localCopiesPath string, sc models.IStorageClaim) files.Filer {
	f := files.NewPath(localCopiesPath).BuildFile(sc.GetName())

	return f
}

func isLocalStorageClaim(sc models.IStorageClaim) bool {
	_, is := sc.(*StorageClaim)

	return is
}

// MoveFileToStorageClaim _
func (client *StorageClient) MoveFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.Progresser) error {
	// TODO: just move file for local storage

	err := client.UploadFileToStorageClaim(ctx, file, sc, p)

	if err != nil {
		return errors.Wrap(err, "Uploading")
	}

	err = file.Remove()

	if err != nil {
		return errors.Wrap(err, "Removing file after uploading to storage claim")
	}

	return nil
}

// UploadFileToStorageClaim _
func (client *StorageClient) UploadFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.Progresser) error {
	// TODO: notify progress
	// TODO: handle ctx cancel

	storageWriter, err := sc.GetWriter()

	if err != nil {
		return errors.Wrap(err, "Building storage claim writer")
	}

	fileReader, err := file.ReadContent()

	if err != nil {
		return errors.Wrap(err, "Building file reader")
	}

	ctxReader := ctxio.NewReader(ctx, fileReader)

	progressReader := progress.NewReader(ctxReader)

	_, err = io.Copy(storageWriter, progressReader)

	storageWriter.Close()
	fileReader.Close()

	if err != nil {
		return errors.Wrap(err, "Failed to upload")
	}

	return nil
}

// DownloadFileFromStorageClaim _
func (client *StorageClient) DownloadFileFromStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.Progresser) error {
	// TODO: notify progress
	// TODO: handle ctx cancel

	fileWriter, err := file.WriteContent()

	if err != nil {
		return errors.Wrap(err, "Building file writer")
	}

	storageReader, err := sc.GetReader()

	if err != nil {
		return errors.Wrap(err, "Building storage claim reader")
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
		return errors.Wrap(err, "Failed to download")
	}

	return nil
}

// BuildStorageClaim _
// func (client *StorageClient) BuildStorageClaim(identity string) (models.IStorageClaim, error) {
// 	claimFile := sc.storagePath.BuildFile(identity)

// 	if claimFile.IsExist() == false {
// 		return nil, ErrStorageClaimMissingFile
// 	}

// 	size, err := claimFile.Size()

// 	if err != nil {
// 		return nil, errors.Wrap(err, "Getting claim file size")
// 	}

// 	return &StorageClaim{
// 		identity: identity,
// 		file:     claimFile,
// 		size:     size,
// 	}, nil
// }
