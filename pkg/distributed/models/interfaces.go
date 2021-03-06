package models

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/wailorman/fftb/pkg/files"
)

// ErrUnknownRequestType _
var ErrUnknownRequestType = errors.New("Unknown request type")

// ErrUnknownOrderType _
var ErrUnknownOrderType = errors.New("Unknown order type")

// ErrUnknownSegmentType _
var ErrUnknownSegmentType = errors.New("Unknown segment type")

// ErrNotImplemented _
var ErrNotImplemented = errors.New("Not implemented")

// ErrUnknownStorageClaimType _
var ErrUnknownStorageClaimType = errors.New("Unknown storage claim type")

// ErrMissingStorageClaim _
var ErrMissingStorageClaim = errors.New("Missing storage claim")

// ErrNotFound _
var ErrNotFound = errors.New("Not found")

// ErrTimeoutReached _
var ErrTimeoutReached = errors.New("Timeout reached")

// ErrFreeSegmentLockTimeout _
var ErrFreeSegmentLockTimeout = errors.New("Free segment lock timeout") // TODO: Subst. with ErrLockTimeoutReached

// ErrLockTimeoutReached _
var ErrLockTimeoutReached = errors.New("Lock timeout reached")

// ErrMissingLockAuthor _
var ErrMissingLockAuthor = errors.New("Missing lock author")

// ErrMissingSegment _
var ErrMissingSegment = errors.New("Missing Segment")

// ErrMissingOrder _
var ErrMissingOrder = errors.New("Missing Order")

// ErrMissingPublisher _
var ErrMissingPublisher = errors.New("Missing publisher")

// ErrMissingPerformer _
var ErrMissingPerformer = errors.New("Missing performer")

// ErrPerformerMismatch _
var ErrPerformerMismatch = errors.New("Performer mismatch")

// ErrCancelled _
var ErrCancelled = errors.New("Cancelled")

// ProgressStep _
type ProgressStep string

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

// SegmentLockDuration _
const SegmentLockDuration = time.Duration(1 * time.Minute)

// // LocalAuthorName _
// var LocalAuthorName = "local"

// // LocalAuthor _
// var LocalAuthor = &Author{name: LocalAuthorName}

// Author _
type Author struct {
	Name string `json:"name"`
}

// GetName _
func (a *Author) GetName() string {
	return a.Name
}

// IsEqual _
func (a *Author) IsEqual(anotherAuthor IAuthor) bool {
	return a.Name == anotherAuthor.GetName()
}

// IContracterRequest _
type IContracterRequest interface {
	GetType() string
	// GetAuthor() IAuthor
	Validate() error
}

// IDealerRequest _
type IDealerRequest interface {
	GetID() string
	GetType() string
	GetAuthor() IAuthor
	Validate() error
}

// IOrder _
type IOrder interface {
	GetID() string
	// SetID(string)
	GetType() string
	// SetType(string)
	// GetSegments() []ISegment
	// GetPayload() (string, error)
	GetPublisher() IAuthor
	// SetPublisher(IAuthor)
	// TODO: use specific type for state
	GetState() string
	// SetState(string)
	MatchPublisher(IAuthor) bool
	GetSegmentIDs() []string
	Validate() error
	GetInputFile() files.Filer
	GetOutputFile() files.Filer
	CalculateProgress([]ISegment) float64
}

// ISegment _
type ISegment interface {
	GetID() string
	GetOrderID() string
	GetType() string
	GetInputStorageClaimIdentity() string
	GetOutputStorageClaimIdentity() string
	// GetStorageClaim() IStorageClaim // should be done by dealer
	// GetPayload() (string, error)
	GetIsLocked() bool
	GetLockedBy() IAuthor
	GetLockedUntil() *time.Time
	// TODO: use specific type for state
	GetState() string
	GetCurrentState() string
	GetPublisher() IAuthor
	GetPerformer() IAuthor
	GetPosition() int
	MatchPublisher(IAuthor) bool
	MatchPerformer(IAuthor) bool
	Lock(performer IAuthor)
	Unlock()
	Validate() error
}

// IContracter _
type IContracter interface {
	// PrepareOrder(req IContracterRequest) (IOrder, error)
	GetAllOrders(ctx context.Context) ([]IOrder, error)
	GetAllSegments(ctx context.Context) ([]ISegment, error)
	GetSegmentsByOrderID(fctx context.Context, orderID string) ([]ISegment, error)
	GetSegmentByID(segmentID string) (ISegment, error)
}

// IDealer _
type IDealer interface {
	IContracterDealer
	IWorkerDealer
}

// IContracterDealer _
type IContracterDealer interface {
	GetOutputStorageClaim(publisher IAuthor, segmentID string) (IStorageClaim, error)
	AllocatePublisherAuthority(name string) (IAuthor, error)
	AllocateSegment(req IDealerRequest) (ISegment, error)
	AllocateInputStorageClaim(publisher IAuthor, id string) (IStorageClaim, error)
	// FindSegmentByID(id string) (ISegment, error)
	NotifyRawUpload(publisher IAuthor, id string, p Progresser) error
	NotifyResultDownload(publisher IAuthor, id string, p Progresser) error
	PublishSegment(publisher IAuthor, id string) error
	CancelSegment(publisher IAuthor, id string) error
	WaitOnSegmentFinished(ctx context.Context, id string) <-chan struct{}
	WaitOnSegmentFailed(ctx context.Context, id string) <-chan error
	GetQueuedSegmentsCount(fctx context.Context, publisher IAuthor) (int, error)
	GetSegmentsByOrderID(fctx context.Context, orderID string) ([]ISegment, error)
	GetSegmentsStatesByOrderID(fctx context.Context, orderID string) (map[string]string, error)
	GetSegmentByID(segmentID string) (ISegment, error)
}

// IWorkerDealer _
type IWorkerDealer interface {
	GetInputStorageClaim(performer IAuthor, segmentID string) (IStorageClaim, error)
	AllocatePerformerAuthority(name string) (IAuthor, error)
	FindFreeSegment(performer IAuthor) (ISegment, error)
	NotifyRawDownload(performer IAuthor, id string, p Progresser) error
	NotifyResultUpload(performer IAuthor, id string, p Progresser) error
	NotifyProcess(performer IAuthor, id string, p Progresser) error
	FinishSegment(performer IAuthor, id string) error
	QuitSegment(performer IAuthor, id string) error
	FailSegment(performer IAuthor, id string, err error) error
	AllocateOutputStorageClaim(performer IAuthor, id string) (IStorageClaim, error)
	WaitOnSegmentCancelled(ctx context.Context, id string) <-chan struct{}
}

// IRegistry _
type IRegistry interface {
	FindOrderByID(id string) (IOrder, error)
	PersistOrder(order IOrder) error
	PickOrderFromQueue(context.Context) (IOrder, error)
	FindSegmentByID(id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	PersistSegment(ISegment) error
	SearchOrder(fctx context.Context, check func(IOrder) bool) (IOrder, error)
	SearchAllOrders(fctx context.Context, check func(IOrder) bool) ([]IOrder, error)
	FindNotLockedSegment(ctx context.Context) (ISegment, error)
	// LockSegmentByID(segmentID string, lockedBy IAuthor) error
	// UnlockSegmentByID(segmentID string) error
	SearchSegment(fctx context.Context, check func(ISegment) bool) (ISegment, error)
	SearchAllSegments(fctx context.Context, check func(ISegment) bool) ([]ISegment, error)
	Persist() error
	Closed() <-chan struct{}
}

// IContracterRegistry _
type IContracterRegistry interface {
	// Persist(ISegment) error
	// FindByID(string) (ISegment, error)
	// Destroy(string) (errors)

	FindOrderByID(id string) (IOrder, error)
	PersistOrder(order IOrder) error
	PickOrderFromQueue(context.Context) (IOrder, error)

	FindSegmentByID(id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	// FindNotLockedSegment(ctx context.Context) (ISegment, error)
	PersistSegment(ISegment) error
	// LockSegmentByID(segmentID string, lockedBy IAuthor) error
	// UnlockSegmentByID(segmentID string) error
	// SearchSegment(fctx context.Context, check func(ISegment) bool) (ISegment, error)
	// SearchAllSegments(fctx context.Context, check func(ISegment) bool) ([]ISegment, error)
	SearchOrder(fctx context.Context, check func(IOrder) bool) (IOrder, error)
	SearchAllOrders(fctx context.Context, check func(IOrder) bool) ([]IOrder, error)
	Persist() error
	Closed() <-chan struct{}
}

// IDealerRegistry _
type IDealerRegistry interface {
	// Persist(ISegment) error
	// FindByID(string) (ISegment, error)
	// Destroy(string) (errors)

	// FindOrderByID(id string) (IOrder, error)
	// PersistOrder(order IOrder) error
	// PickOrderFromQueue(context.Context) (IOrder, error)

	FindSegmentByID(id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	FindNotLockedSegment(ctx context.Context) (ISegment, error)
	PersistSegment(ISegment) error
	// LockSegmentByID(segmentID string, lockedBy IAuthor) error
	// UnlockSegmentByID(segmentID string) error
	SearchSegment(fctx context.Context, check func(ISegment) bool) (ISegment, error)
	SearchAllSegments(fctx context.Context, check func(ISegment) bool) ([]ISegment, error)
	Persist() error
	Closed() <-chan struct{}
}

// IStorageController _
type IStorageController interface {
	AllocateStorageClaim(name string) (IStorageClaim, error)
	PurgeStorageClaim(claim IStorageClaim) error
	BuildStorageClaim(name string) (IStorageClaim, error)
}

// IStorageClaim _
type IStorageClaim interface {
	GetID() string
	GetName() string
	GetSize() (int, error)
	GetWriter() (io.WriteCloser, error)
	GetReader() (io.ReadCloser, error)
}

// Progresser _
type Progresser interface {
	Step() ProgressStep
	Percent() float64
}

// Subscriber _
type Subscriber interface {
	GetOutput() chan Progresser
	Unsubscribe()
}

// PublishSubscriber _
type PublishSubscriber interface {
	Subscribe() Subscriber
	Publish(Progresser)
}

// IAuthor _
type IAuthor interface {
	GetName() string
	IsEqual(IAuthor) bool
}
