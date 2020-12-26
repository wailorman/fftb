package local

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// Dealer _
type Dealer struct {
	storageController models.IStorageController
	registry          models.IRegistry
}

// NewDealer _
func NewDealer(sc models.IStorageController, r models.IRegistry) *Dealer {
	return &Dealer{
		storageController: sc,
		registry:          r,
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
	claimIdentity := fmt.Sprintf("%s/%s/%s", convertReq.OrderIdentity, convertReq.Identity, id)

	claim, err := d.storageController.AllocateStorageClaim(claimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating storage claim for task")
	}

	convertSegment := &models.ConvertSegment{
		Identity:             convertReq.Identity,
		OrderIdentity:        convertReq.OrderIdentity,
		StorageClaimIdentity: claim.GetID(),
		Params:               convertReq.Params,

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
func (d *Dealer) FindFreeSegment() (models.ISegment, error) {
	panic(models.ErrNotImplemented)
	// return nil, nil
}

// GetStorageClaim _
func (d *Dealer) GetStorageClaim(segment models.ISegment) (models.IStorageClaim, error) {
	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.StorageClaimIdentity)

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
