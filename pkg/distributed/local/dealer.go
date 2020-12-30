package local

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// LocalAuthor _
const LocalAuthor = "local"

// Dealer _
type Dealer struct {
	storageController models.IStorageController
	registry          models.IRegistry
	freeSegmentLock   trylock.TryLocker
}

// NewDealer _
func NewDealer(sc models.IStorageController, r models.IRegistry) *Dealer {
	return &Dealer{
		storageController: sc,
		registry:          r,
		freeSegmentLock:   trylock.New(),
	}
}

// AllocateSegment _
func (d *Dealer) AllocateSegment(req models.IDealerRequest) (models.ISegment, error) {
	convertReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return nil, models.ErrUnknownRequestType
	}

	// fmt.Printf("convertReq.Identity: %#v\n", convertReq.Identity)
	// id := fmt.Sprintf("%s/%s", convertReq.Identity, uuid.New().String())
	// claimIdentity := fmt.Sprintf("%s.%s", id, convertReq.Params.Muxer)
	// claimIdentity := fmt.Sprintf("%s.%s", id, "mp4")
	id := uuid.New().String()
	inputClaimIdentity := fmt.Sprintf("%s/%s/%s_input", convertReq.OrderIdentity, convertReq.Identity, id)
	outputClaimIdentity := fmt.Sprintf("%s/%s/%s_output", convertReq.OrderIdentity, convertReq.Identity, id)

	inputClaim, err := d.storageController.AllocateStorageClaim(inputClaimIdentity)
	outputClaim, err := d.storageController.AllocateStorageClaim(outputClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating storage claim for task")
	}

	convertSegment := &models.ConvertSegment{
		Identity:                   convertReq.Identity,
		OrderIdentity:              convertReq.OrderIdentity,
		InputStorageClaimIdentity:  inputClaim.GetID(),
		OutputStorageClaimIdentity: outputClaim.GetID(),
		Params:                     convertReq.Params,

		// Muxer:      convertReq.Muxer,
		// VideoCodec: convertReq.VideoCodec,
		// // HWAccel:          convertReq.HWAccel,
		// // VideoBitRate:     convertReq.VideoBitRate,
		// VideoQuality: convertReq.VideoQuality,
		// // Preset:           convertReq.Preset,
		// // Scale:            convertReq.Scale,
		// // KeyframeInterval: convertReq.KeyframeInterval,
	}

	// set state to @prepared

	return convertSegment, nil
}

// FindFreeSegment _
func (d *Dealer) FindFreeSegment(author string) (models.ISegment, error) {
	locked := d.freeSegmentLock.TryLockTimeout(LockSegmentTimeout)

	if !locked {
		return nil, models.ErrFreeSegmentLockTimeout
	}

	defer d.freeSegmentLock.Unlock()

	freeSegment, err := d.registry.FindNotLockedSegment()

	if err != nil {
		return nil, errors.Wrap(err, "Looking for free segment")
	}

	err = d.registry.LockSegmentByID(freeSegment.GetID(), author)

	if err != nil {
		return nil, errors.Wrap(err, "Locking free segment")
	}

	return freeSegment, nil
}

// GetInputStorageClaim _
func (d *Dealer) GetInputStorageClaim(segment models.ISegment) (models.IStorageClaim, error) {
	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.InputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

// GetOutputStorageClaim _
func (d *Dealer) GetOutputStorageClaim(segment models.ISegment) (models.IStorageClaim, error) {
	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.OutputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

// CancelSegment _
func (d *Dealer) CancelSegment(models.ISegment) error {
	panic(models.ErrNotImplemented)
}

// FindSegmentByID _
func (d *Dealer) FindSegmentByID(id string) (models.ISegment, error) {
	panic(models.ErrNotImplemented)
}

// NotifyRawUpload _
func (d *Dealer) NotifyRawUpload(progresser models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyResultDownload _
func (d *Dealer) NotifyResultDownload(progresser models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// PublishSegment _
func (d *Dealer) PublishSegment(segment models.ISegment) error {
	return d.registry.PersistSegment(segment)
	// panic(models.ErrNotImplemented)
}

// Subscription _
func (d *Dealer) Subscription(segment models.ISegment) (models.Subscriber, error) {
	panic(models.ErrNotImplemented)
}

// FinishSegment _
func (d *Dealer) FinishSegment(models.ISegment, models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyProcess _
func (d *Dealer) NotifyProcess(models.ISegment, models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyRawDownload _
func (d *Dealer) NotifyRawDownload(models.ISegment, models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyResultUpload _
func (d *Dealer) NotifyResultUpload(models.ISegment, models.Progresser) error {
	panic(models.ErrNotImplemented)
}
