package s3

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buildStorageController() *StorageControl {
	return NewStorageControl()
}

func buildStorageClient() *StorageClient {
	return NewStorageClient(".tmp/s3")
}

func Test__Full(t *testing.T) {
	storageController := buildStorageController()
	storageClient := buildStorageClient()
	claimName := "fftb_s3_test_" + uuid.New().String()
	content := "Lorem ipsum"

	allocatedStorageClaim, err := storageController.AllocateStorageClaim(context.Background(), claimName)

	if !assert.NoError(t, err, "Failed to allocate storage claim") {
		return
	}

	assert.NotEqual(t, "", allocatedStorageClaim.GetURL(), "allocated storage claim url")

	clientStorageClaim, err := storageClient.BuildStorageClaimByURL(allocatedStorageClaim.GetURL())

	if !assert.NoError(t, err, "Failed to build storage claim on client") {
		return
	}

	// writer, err := clientStorageClaim.GetWriter()

	// if !assert.NoError(t, err, "Failed to build storage claim writer") {
	// 	return
	// }

	// contentUploadingReader := strings.NewReader(content)

	// _, err = io.Copy(writer, contentUploadingReader)

	// writer.Close()

	// if !assert.NoError(t, err, "Failed to write content to storage claim") {
	// 	return
	// }

	contentReader := strings.NewReader(content)
	err = clientStorageClaim.WriteFrom(contentReader)

	if !assert.NoError(t, err, "Failed to write content to storage claim") {
		return
	}

	downloadingStorageClaimURL, err := storageController.BuildStorageClaim(allocatedStorageClaim.GetID())

	if !assert.NoError(t, err, "Failed to build downloading URL") {
		return
	}

	downloadingStorageClaim, err := storageClient.BuildStorageClaimByURL(downloadingStorageClaimURL.GetURL())

	if !assert.NoError(t, err, "Failed to build downloading storage claim") {
		return
	}

	// contentDownloadingReader, err := downloadingStorageClaim.GetReader()

	// if !assert.NoError(t, err, "Failed to build downloading reader") {
	// 	return
	// }

	// downloadedContentRaw, err := io.ReadAll(contentDownloadingReader)

	// if !assert.NoError(t, err, "Failed to download content") {
	// 	return
	// }

	// downloadedContent := string(downloadedContentRaw)
	buf := bytes.NewBufferString("")

	err = downloadingStorageClaim.ReadTo(buf)

	if !assert.NoError(t, err, "Failed to download content") {
		return
	}

	assert.Equal(t, content, buf.String(), "Wrong content")
}
