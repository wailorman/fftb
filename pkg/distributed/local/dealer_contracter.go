package local

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// AllocateSegment _
func (d *Dealer) AllocateSegment(req models.IDealerRequest) (models.ISegment, error) {
	convertReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return nil, models.ErrUnknownRequestType
	}

	// TODO: check id is free

	convertSegment := &models.ConvertSegment{
		Identity:      convertReq.Identity,
		OrderIdentity: convertReq.OrderIdentity,
		Params:        convertReq.Params,
		Muxer:         convertReq.Muxer,
		State:         models.SegmentPreparedState,
		Publisher:     req.GetAuthor(),
	}

	// TODO: persist

	err := d.registry.PersistSegment(convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting segment")
	}

	return convertSegment, nil
}

// GetOutputStorageClaim _
func (d *Dealer) GetOutputStorageClaim(publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	if convertSegment.OutputStorageClaimIdentity == "" {
		return nil, errors.Wrap(models.ErrMissingStorageClaim, "Getting output storage claim identity")
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.OutputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

// AllocateInputStorageClaim _
func (d *Dealer) AllocateInputStorageClaim(publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	iClaimID := fmt.Sprintf("%s/%s/input_%s", segment.GetOrderID(), segment.GetID(), uuid.New().String())
	iClaim, err := d.storageController.AllocateStorageClaim(iClaimID)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating input storage claim")
	}

	convertSegment.InputStorageClaimIdentity = iClaimID

	err = d.registry.PersistSegment(convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting input claim identity")
	}

	return iClaim, nil
}

// CancelSegment _
func (d *Dealer) CancelSegment(publisher models.IAuthor, segmentID string) error {
	panic(models.ErrNotImplemented)
}

// NotifyRawUpload _
func (d *Dealer) NotifyRawUpload(publisher models.IAuthor, segmentID string, p models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyResultDownload _
func (d *Dealer) NotifyResultDownload(publisher models.IAuthor, segmentID string, p models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// PublishSegment _
func (d *Dealer) PublishSegment(publisher models.IAuthor, segmentID string) error {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	convertSegment.State = models.SegmentPublishedState

	return d.registry.PersistSegment(convertSegment)
}

// AllocatePublisherAuthority _
func (d *Dealer) AllocatePublisherAuthority(name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/publishers/%s", name)

	return &models.Author{Name: authorName}, nil
}
