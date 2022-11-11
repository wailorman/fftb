package models

import (
	"context"
	"io"

	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ProgressStep _
type ProgressStep string

type StorageClaimPurpose string

// UploadingInputStep _
const UploadingInputStep ProgressStep = "uploading_input"

// DownloadingInputStep _
const DownloadingInputStep ProgressStep = "downloading_input"

// ProcessingStep _
const ProcessingStep ProgressStep = "processing"

// UploadingOutputStep _
const UploadingOutputStep ProgressStep = "uploading_output"

// DownloadingOutputStep _
const DownloadingOutputStep ProgressStep = "downloading_output"

const NoneStorageClaimPurpose StorageClaimPurpose = "none"
const ConvertInputStorageClaimPurpose StorageClaimPurpose = "convert_input"
const ConvertOutputStorageClaimPurpose StorageClaimPurpose = "convert_output"

// LocalAuthorName _
var LocalAuthorName = "local"

// LocalAuthor _
var LocalAuthor = &Author{Name: LocalAuthorName}

// Author _
type Author struct {
	Name string `json:"name"`
}

// GetName _
func (a *Author) GetName() string {
	return a.Name
}

// ISegment _
type ISegment interface {
	GetID() string
	GetType() pb.SegmentType
	GetConvertSegmentParams() convert.Params
}

// IDealer _
type IDealer interface {
	IWorkerDealer
}

type StorageClaimRequest struct {
	SegmentID string
	Purpose   StorageClaimPurpose
	Name      string
}

// IWorkerDealer _
type IWorkerDealer interface {
	AllocatePerformerAuthority(ctx context.Context, name string) (IAuthor, error)

	// TODO: add search criteria
	FindFreeSegment(ctx context.Context, performer IAuthor) (ISegment, error)

	NotifyRawDownload(ctx context.Context, performer IAuthor, segmentID string, p IProgress) error
	NotifyResultUpload(ctx context.Context, performer IAuthor, segmentID string, p IProgress) error
	NotifyProcess(ctx context.Context, performer IAuthor, segmentID string, p IProgress) error

	FinishSegment(ctx context.Context, performer IAuthor, id string) error
	QuitSegment(ctx context.Context, performer IAuthor, id string) error
	FailSegment(ctx context.Context, performer IAuthor, id string, err error) error

	GetAllInputStorageClaims(ctx context.Context, performer IAuthor, req StorageClaimRequest) ([]IStorageClaim, error)
	AllocateOutputStorageClaim(ctx context.Context, performer IAuthor, req StorageClaimRequest) (IStorageClaim, error)
}

// IStorageClaim _
type IStorageClaim interface {
	GetID() string
	GetName() string
	GetURL() string
	GetSize() int
	GetType() string
	// GetWriter() (io.WriteCloser, error)
	// GetReader() (io.ReadCloser, error)
	// WriteFrom(io.ReadCloser) error
	// ReadTo(io.WriteCloser) error
	WriteFrom(io.Reader) error
	ReadTo(io.Writer) error
}

// IStorageClient _
type IStorageClient interface {
	// TODO: ctx
	BuildStorageClaimByURL(url string) (IStorageClaim, error)
	RemoveLocalCopy(ctx context.Context, sc IStorageClaim) error
	MakeLocalCopy(ctx context.Context, sc IStorageClaim, p chan IProgress) (files.Filer, error)
	MoveFileToStorageClaim(ctx context.Context, file files.Filer, sc IStorageClaim, p chan IProgress) error
}

// IProgress _
type IProgress interface {
	Step() ProgressStep
	Percent() float64
}

// IAuthor _
type IAuthor interface {
	GetName() string
}

// EmptyAuthor TODO: REMOVE
var EmptyAuthor = &Author{}
