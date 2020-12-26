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

// IContractDealer _
type IContractDealer interface {
	AllocateSegment(req IDealerRequest) (ISegment, error)
	FindSegmentByID(id string) (ISegment, error)
	GetStorageClaim(ISegment) (IStorageClaim, error)
	NotifyRawUpload(Progresser) error
	NotifyResultDownload(Progresser) error
	PublishSegment(ISegment) error
	CancelSegment(ISegment) error
	Subscription(ISegment) (Subscriber, error)
}

// IWorkDealer _
type IWorkDealer interface {
	FindFreeSegment() (ISegment, error)
	GetStorageClaim(ISegment) IStorageClaim
	NotifyRawDownload(Progresser) error
	NotifyResultUpload(Progresser) error
	NotifyProcess(Progresser) error
	FinishSegment(Progresser) error
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
	GetStorageClaimIdentity() string
	// GetStorageClaim() IStorageClaim // should be done by dealer
	GetPayload() (string, error)
}

// IRegistry _
type IRegistry interface {
	// Persist(ISegment) error
	// FindByID(string) (ISegment, error)
	// Destroy(string) (errors)

	FindSegmentByID(id string) (ISegment, error)
	FindSegmentsByOrderID(orderID string) ([]ISegment, error)
	PersistSegment(segment ISegment) error
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
