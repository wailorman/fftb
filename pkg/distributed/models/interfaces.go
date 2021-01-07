package models

import (
	"errors"
	"io"
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

// // ErrNoFreeSegments _
// var ErrNoFreeSegments = errors.New("No free segments")

// ErrFreeSegmentLockTimeout _
var ErrFreeSegmentLockTimeout = errors.New("Free segment lock timeout")

// ErrMissingLockAuthor _
var ErrMissingLockAuthor = errors.New("Missing lock author")

// ErrMissingSegment _
var ErrMissingSegment = errors.New("Missing Segment")

// ErrMissingOrder _
var ErrMissingOrder = errors.New("Missing Order")

// IContracter _
type IContracter interface {
	PrepareOrder(req IContracterRequest) (IOrder, error)
}

// IContracterRequest _
type IContracterRequest interface {
	GetType() string
}

// IOrder _
type IOrder interface {
	GetID() string
	GetType() string
	GetSegments() []ISegment
	GetPayload() (string, error)
	Failed(error)
}

// InputOutputStorageClaimer _
type InputOutputStorageClaimer interface {
	GetInputStorageClaim(ISegment) (IStorageClaim, error)
	GetOutputStorageClaim(ISegment) (IStorageClaim, error)
}

// IContractDealer _
type IContractDealer interface {
	InputOutputStorageClaimer

	AllocateSegment(req IDealerRequest) (ISegment, error)
	FindSegmentByID(id string) (ISegment, error)
	NotifyRawUpload(Progresser) error
	NotifyResultDownload(Progresser) error
	PublishSegment(ISegment) error
	CancelSegment(ISegment) error
	Subscription(ISegment) (Subscriber, error)
	AllocateInputStorageClaim(segment ISegment) (IStorageClaim, error)
}

// IWorkDealer _
type IWorkDealer interface {
	InputOutputStorageClaimer

	FindFreeSegment(author string) (ISegment, error)
	NotifyRawDownload(ISegment, Progresser) error
	NotifyResultUpload(ISegment, Progresser) error
	NotifyProcess(ISegment, Progresser) error
	FinishSegment(ISegment) error
	AllocateOutputStorageClaim(segment ISegment) (IStorageClaim, error)
}

// IDealerRequest _
type IDealerRequest interface {
	GetID() string
	GetType() string
}

// ISegment _
type ISegment interface {
	GetID() string
	GetOrderID() string
	GetType() string
	GetInputStorageClaimIdentity() string
	GetOutputStorageClaimIdentity() string
	// GetStorageClaim() IStorageClaim // should be done by dealer
	GetPayload() (string, error)
	GetIsLocked() bool
	GetLockedBy() string
	// TODO: use specific type for segment state
	GetState() string
}

// IRegistry _
type IRegistry interface {
	// Persist(ISegment) error
	// FindByID(string) (ISegment, error)
	// Destroy(string) (errors)

	FindSegmentByID(id string) (ISegment, error)
	FindSegmentsByOrderID(orderID string) ([]ISegment, error)
	FindNotLockedSegment() (ISegment, error)
	PersistSegment(segment ISegment) error
	LockSegmentByID(segmentID string, lockedBy string) error
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
	Step() string
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
