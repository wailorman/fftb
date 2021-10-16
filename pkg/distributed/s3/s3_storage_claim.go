// TODO: move under storage package
package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// Timeout _
var Timeout = time.Duration(60 * time.Second)

// DefaultContentType _
const DefaultContentType = "application/octet-stream"

// StorageClaimTypeS3 is remote storage claim type, working on S3.
// Client works with http
const StorageClaimTypeS3 = "s3"

// StorageClaim _
type StorageClaim struct {
	id   string
	url  string
	size int
}

// GetID _
func (s *StorageClaim) GetID() string {
	return s.id
}

// GetURL _
func (s *StorageClaim) GetURL() string {
	return s.url
}

// GetName _
func (s *StorageClaim) GetName() string {
	return s.id
}

// GetSize _
func (s *StorageClaim) GetSize() int {
	return s.size
}

// GetType _
func (s *StorageClaim) GetType() string {
	return StorageClaimTypeS3
}

// SplittedIO _
type SplittedIO struct {
	writer io.Writer
	reader io.Reader
	closer io.Closer
}

// NewSplittedIO _
func NewSplittedIO(writer io.Writer, reader io.Reader, closer io.Closer) *SplittedIO {
	return &SplittedIO{
		writer: writer,
		reader: reader,
		closer: closer,
	}
}

// Read _
func (sio *SplittedIO) Read(p []byte) (n int, err error) {
	if sio.reader == nil {
		return 0, nil
	}

	return sio.reader.Read(p)
}

// Write _
func (sio *SplittedIO) Write(p []byte) (n int, err error) {
	if sio.writer == nil {
		return 0, nil
	}

	return sio.writer.Write(p)
}

// Close _
func (sio *SplittedIO) Close() error {
	if sio.closer == nil {
		return nil
	}

	return sio.closer.Close()
}

// // NewSplittedWriteCloser _
// func NewSplittedWriteCloser(writer io.Writer, closer io.Closer) *SplittedWriteCloser {
// 	return &SplittedWriteCloser{
// 		writer: writer,
// 		closer: closer,
// 	}
// }

// // SplittedWriteCloser _
// type SplittedWriteCloser struct {
// 	writer io.Writer
// 	closer io.Closer
// }

// // Write _
// func (wc *SplittedWriteCloser) Write(p []byte) (n int, err error) {
// 	return wc.writer.Write(p)
// }

// // Close _
// func (wc *SplittedWriteCloser) Close() error {
// 	return wc.closer.Close()
// }

// // NewSplittedReadCloser _
// func NewSplittedReadCloser(writer io.Reader, closer io.Closer) *SplittedReadCloser {
// 	return &SplittedReadCloser{
// 		writer: writer,
// 		closer: closer,
// 	}
// }

// // SplittedReadCloser _
// type SplittedReadCloser struct {
// 	writer io.Reader
// 	closer io.Closer
// }

// // Read _
// func (wc *SplittedReadCloser) Read(p []byte) (n int, err error) {
// 	return wc.writer.Read(p)
// }

// // Close _
// func (wc *SplittedReadCloser) Close() error {
// 	return wc.closer.Close()
// }

// CombinedCloser _
type CombinedCloser struct {
	closers []io.Closer
}

// NewCombinedCloser _
func NewCombinedCloser(closers ...io.Closer) *CombinedCloser {
	if closers == nil {
		return &CombinedCloser{closers: make([]io.Closer, 0)}
	}

	return &CombinedCloser{closers: closers}
}

// Close _
func (cc *CombinedCloser) Close() error {
	var lastErr error

	for _, closer := range cc.closers {
		if err := closer.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// WriteFrom _
func (s *StorageClaim) WriteFrom(reader io.Reader) error {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	// u, err := url.Parse(s.url)

	// if err != nil {
	// 	return errors.Wrap(err, "Failed to parse claim url")
	// }

	// host := u.Host

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return errors.Wrap(err, "Failed to buffer reader")
	}

	req, err := http.NewRequestWithContext(context.TODO(), "PUT", s.url, &buf)
	defer req.Body.Close()
	req.Header.Set("Content-Type", DefaultContentType)
	// req.Header.Set("Host", host)
	if err != nil {
		return errors.Wrap(err, "Building request")
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "Performing upload request")
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.Wrapf(wrapHTTPError(res), "Failed upload request (url: `%s`)", s.url)
	}

	return nil
}

// ReadTo _
func (s *StorageClaim) ReadTo(writer io.Writer) error {
	if closer, ok := writer.(io.Closer); ok {
		defer closer.Close()
	}

	res, err := http.Get(s.url)

	if err != nil {
		return errors.Wrap(err, "Performing download request")
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.Wrapf(wrapHTTPError(res), "Failed download request (url: `%s`)", s.url)
	}

	_, err = io.Copy(writer, res.Body)

	if err != nil {
		return errors.Wrap(err, "Writing response to writer")
	}

	return nil
}

func wrapHTTPError(resp *http.Response) error {
	kylobytes := 1024
	if resp.ContentLength > 0 && resp.ContentLength < int64(2*kylobytes) {
		rawBody, _ := io.ReadAll(resp.Body)
		return errors.Wrapf(models.ErrUnknown, "%s `%s`", resp.Status, string(rawBody))
	}

	return errors.Wrapf(models.ErrUnknown, "HTTP %s", resp.Status)
}
