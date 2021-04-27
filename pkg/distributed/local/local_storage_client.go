package local

import (
	"context"
	"strings"

	"github.com/pkg/errors"
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

// BuildStorageClaimByURL _
func (client *StorageClient) BuildStorageClaimByURL(url string) (models.IStorageClaim, error) {
	if strings.Index(url, "file://") != 0 {
		return nil, errors.Wrapf(models.ErrUnknownType, "Unknown storage claim type passed in url `%s`", url)
	}

	claimPath := strings.Replace(url, "file://", "", 1)
	file := files.NewFile(claimPath)
	size, err := file.Size()

	if err != nil {
		return nil, errors.Wrap(err, "Getting file size")
	}

	claim := &StorageClaim{
		identity: file.Name(),
		file:     file,
		size:     size,
	}

	return claim, nil
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

	// storageWriter, err := sc.GetWriter()

	// if err != nil {
	// 	return errors.Wrap(err, "Building storage claim writer")
	// }

	fileReader, err := file.ReadContent()

	if err != nil {
		return errors.Wrap(err, "Building file reader")
	}

	err = sc.WriteFrom(fileReader)

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

	err = sc.ReadTo(fileWriter)

	if err != nil {
		return errors.Wrap(err, "Writing storage claim content to file")
	}

	return nil
}
