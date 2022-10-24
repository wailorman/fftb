package s3

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
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
func (client *StorageClient) BuildStorageClaimByURL(sURL string) (models.IStorageClaim, error) {
	if strings.Index(sURL, "http") != 0 {
		return nil, errors.Wrapf(models.ErrUnknownType, "Unknown storage claim type passed in url `%s`", sURL)
	}

	sizeReq, err := http.NewRequest("HEAD", sURL, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Building HEAD request for calculating size")
	}

	httpClient := http.Client{Timeout: Timeout}

	sizeRes, err := httpClient.Do(sizeReq)

	if err != nil {
		return nil, errors.Wrap(err, "Performing HEAD request for calculating size")
	}

	sizeStr := sizeRes.Header.Get("Content-Length")
	size, _ := strconv.Atoi(sizeStr)

	u, err := url.Parse(sURL)

	if err != nil {
		return nil, errors.Wrap(err, "Parsing url")
	}

	// urlParts := strings.Split(url, "/")
	id := u.Path

	claim := &StorageClaim{
		id:   id,
		size: size,
		url:  sURL,
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
func (client *StorageClient) MakeLocalCopy(ctx context.Context, sc models.IStorageClaim, p chan models.IProgress) (files.Filer, error) {
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

// MoveFileToStorageClaim _
func (client *StorageClient) MoveFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.IProgress) error {
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
func (client *StorageClient) UploadFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.IProgress) error {
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
func (client *StorageClient) DownloadFileFromStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.IProgress) error {
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

// BuildStorageClaim _
// func (client *StorageClient) BuildStorageClaim(identity string) (models.IStorageClaim, error) {
// 	claimFile := sc.storagePath.BuildFile(identity)

// 	if claimFile.IsExist() == false {
// 		return nil, ErrNotFound
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
