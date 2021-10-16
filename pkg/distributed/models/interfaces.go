package models

import (
	"context"
	"io"
	"time"

	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/files"
)

// // ErrCancelled _
// var ErrCancelled = errors.New("Cancelled")

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

// MaxRetriesCount _
const MaxRetriesCount = 3

// NextRetryOffset _
const NextRetryOffset = time.Duration(30 * time.Second)

// // LocalAuthorName _
// var LocalAuthorName = "local"

// // LocalAuthor _
// var LocalAuthor = &Author{name: LocalAuthorName}

// CancellationReasonFailed _
const CancellationReasonFailed = "failed"

// CancellationReasonByUser _
const CancellationReasonByUser = "by_user"

// CancellationReasonOrderCancelled _
const CancellationReasonOrderCancelled = "order_cancelled"

// CancellationReasonNotAccepted _
const CancellationReasonNotAccepted = "not_accepted"

// Author _
type Author struct {
	Name         string `json:"name"`
	AuthorityKey string `json:"authority_key,omitempty"`
	SessionKey   string `json:"session_key,omitempty"`
}

// GetName _
func (a *Author) GetName() string {
	return a.Name
}

// GetAuthorityKey _
func (a *Author) GetAuthorityKey() string {
	return a.AuthorityKey
}

// GetSessionKey _
func (a *Author) GetSessionKey() string {
	return a.SessionKey
}

// SetAuthorityKey _
func (a *Author) SetAuthorityKey(key string) {
	a.AuthorityKey = key
}

// SetSessionKey _
func (a *Author) SetSessionKey(key string) {
	a.SessionKey = key
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
	// GetAuthor() IAuthor
	// TODO: validate segment & order id format
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
	Validate() error
	GetInputFile() files.Filer
	GetOutputFile() files.Filer
	CalculateProgress([]ISegment) float64
	GetRetriesCount() int
	GetRetryAt() *time.Time
	GetCanRetry() bool
	GetCanPublish() bool
	GetCanConcat(segments []ISegment) bool

	setLastError(err error)
	incrementRetriesCount()
	cancel(reason string)
}

// SegmentCanceller _
type SegmentCanceller interface {
	CancelSegment(segment ISegment, reason string) error
}

// IOrderMutator _
type IOrderMutator interface {
	CancelOrder(segmentMutator SegmentCanceller, order IOrder, segments []ISegment, reason string) error
	FailOrder(segmentMutator SegmentCanceller, order IOrder, segments []ISegment, err error) error
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
	Validate() error
	GetRetriesCount() int
	GetRetryAt() *time.Time
	GetCanRetry() bool
	GetCanPerform() bool

	lock(performer IAuthor)
	unlock()
	incrementRetriesCount()
	setLastError(err error)
	cancel(reason string)
	publish()
	finish()
}

// ISegmentMutator _
type ISegmentMutator interface {
	// UnlockSegment
	PublishSegment(segment ISegment) error
	FinishSegment(segment ISegment) error
	// RepublishSegment(segment ISegment) error
	CancelSegment(segment ISegment, reason string) error
	FailSegment(segment ISegment, err error) error
	LockSegment(segment ISegment, performer IAuthor) error
	UnlockSegment(segment ISegment) error
}

// IContracter _
type IContracter interface {
	GetOrderByID(ctx context.Context, id string) (IOrder, error)
	GetAllOrders(ctx context.Context) ([]IOrder, error)
	SearchAllOrders(ctx context.Context, search IOrderSearchCriteria) ([]IOrder, error)
	GetAllSegments(ctx context.Context) ([]ISegment, error)
	SearchAllSegments(ctx context.Context, search ISegmentSearchCriteria) ([]ISegment, error)
	GetSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	SearchSegmentsByOrderID(ctx context.Context, orderID string, search ISegmentSearchCriteria) ([]ISegment, error)
	GetSegmentByID(ctx context.Context, segmentID string) (ISegment, error)
	CancelOrderByID(ctx context.Context, orderID string, reason string) error
}

// IDealer _
type IDealer interface {
	IContracterDealer
	IWorkerDealer
}

// IContracterDealer _
type IContracterDealer interface {
	AllocatePublisherAuthority(ctx context.Context, name string) (IAuthor, error)
	AllocateSegment(ctx context.Context, publisher IAuthor, req IDealerRequest) (ISegment, error)

	GetOutputStorageClaim(ctx context.Context, publisher IAuthor, segmentID string) (IStorageClaim, error)
	AllocateInputStorageClaim(ctx context.Context, publisher IAuthor, id string) (IStorageClaim, error)

	GetQueuedSegmentsCount(ctx context.Context, publisher IAuthor) (int, error)
	GetSegmentsByOrderID(ctx context.Context, publisher IAuthor, orderID string, search ISegmentSearchCriteria) ([]ISegment, error)
	GetSegmentByID(ctx context.Context, publisher IAuthor, segmentID string) (ISegment, error)

	NotifyRawUpload(ctx context.Context, publisher IAuthor, id string, p Progresser) error
	NotifyResultDownload(ctx context.Context, publisher IAuthor, id string, p Progresser) error

	PublishSegment(ctx context.Context, publisher IAuthor, id string) error
	RepublishSegment(ctx context.Context, publisher IAuthor, id string) error
	CancelSegment(ctx context.Context, publisher IAuthor, id string, reason string) error
	AcceptSegment(ctx context.Context, publisher IAuthor, id string) error

	// TODO: make local func
	ObserveSegments(ctx context.Context, wg chwg.WaitGrouper)
}

// IWorkerDealer _
type IWorkerDealer interface {
	AllocatePerformerAuthority(ctx context.Context, name string) (IAuthor, error)

	// TODO: add search criteria
	FindFreeSegment(ctx context.Context, performer IAuthor) (ISegment, error)

	NotifyRawDownload(ctx context.Context, performer IAuthor, id string, p Progresser) error
	NotifyResultUpload(ctx context.Context, performer IAuthor, id string, p Progresser) error
	NotifyProcess(ctx context.Context, performer IAuthor, id string, p Progresser) error

	FinishSegment(ctx context.Context, performer IAuthor, id string) error
	QuitSegment(ctx context.Context, performer IAuthor, id string) error
	FailSegment(ctx context.Context, performer IAuthor, id string, err error) error

	GetInputStorageClaim(ctx context.Context, performer IAuthor, segmentID string) (IStorageClaim, error)
	AllocateOutputStorageClaim(ctx context.Context, performer IAuthor, id string) (IStorageClaim, error)
}

// IRegistry _
type IRegistry interface {
	FindOrderByID(ctx context.Context, id string) (IOrder, error)
	PersistOrder(ctx context.Context, order IOrder) error
	FindSegmentByID(ctx context.Context, id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	PersistSegment(ctx context.Context, segment ISegment) error
	// TODO: remove. Leave only searchAll
	SearchOrder(ctx context.Context, check func(IOrder) bool) (IOrder, error)
	SearchAllOrders(ctx context.Context, check func(IOrder) bool) ([]IOrder, error)
	// TODO: remove. Leave only searchAll
	SearchSegment(ctx context.Context, check func(ISegment) bool) (ISegment, error)
	SearchAllSegments(ctx context.Context, check func(ISegment) bool) ([]ISegment, error)
	Persist() error
	Closed() <-chan struct{}
}

// IContracterRegistry _
type IContracterRegistry interface {
	FindOrderByID(ctx context.Context, id string) (IOrder, error)
	PersistOrder(ctx context.Context, order IOrder) error

	FindSegmentByID(ctx context.Context, id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	PersistSegment(ctx context.Context, segment ISegment) error
	SearchOrder(ctx context.Context, check func(IOrder) bool) (IOrder, error)
	SearchAllOrders(ctx context.Context, check func(IOrder) bool) ([]IOrder, error)
	Persist() error
	Closed() <-chan struct{}
}

// IDealerRegistry _
type IDealerRegistry interface {
	FindSegmentByID(ctx context.Context, id string) (ISegment, error)
	FindSegmentsByOrderID(ctx context.Context, orderID string) ([]ISegment, error)
	PersistSegment(ctx context.Context, segment ISegment) error
	SearchSegment(ctx context.Context, check func(ISegment) bool) (ISegment, error)
	SearchAllSegments(ctx context.Context, check func(ISegment) bool) ([]ISegment, error)
	Persist() error
	Closed() <-chan struct{}
}

// IStorageController _
type IStorageController interface {
	AllocateStorageClaim(ctx context.Context, name string) (IStorageClaim, error)
	// TODO: receive string (identity) instead of claim
	PurgeStorageClaim(ctx context.Context, claim IStorageClaim) error
	// TODO: ctx
	BuildStorageClaim(name string) (IStorageClaim, error)
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
	MakeLocalCopy(ctx context.Context, sc IStorageClaim, p chan Progresser) (files.Filer, error)
	MoveFileToStorageClaim(ctx context.Context, file files.Filer, sc IStorageClaim, p chan Progresser) error
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
	GetAuthorityKey() string
	SetAuthorityKey(key string)
	GetSessionKey() string
	SetSessionKey(key string)
	IsEqual(IAuthor) bool
}

// TypeStub _
type TypeStub struct {
	Type string `json:"type"`
}
