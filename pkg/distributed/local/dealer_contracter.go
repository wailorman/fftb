package local

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// AllocateSegment _
func (d *Dealer) AllocateSegment(req models.IDealerRequest) (models.ISegment, error) {
	if validationErr := req.Validate(); validationErr != nil {
		return nil, validationErr
	}

	convertReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return nil, models.ErrUnknownRequestType
	}

	d.logger.WithField(dlog.KeyOrderID, convertReq.OrderIdentity).
		WithField(dlog.KeySegmentID, convertReq.Identity).
		Info("Allocating segment")

	// TODO: check id is free

	convertSegment := &models.ConvertSegment{
		Type:          models.ConvertV1Type,
		Identity:      convertReq.Identity,
		OrderIdentity: convertReq.OrderIdentity,
		Params:        convertReq.Params,
		Muxer:         convertReq.Muxer,
		Position:      convertReq.Position,
		State:         models.SegmentStatePrepared,
		Publisher:     req.GetAuthor(),
	}

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

	d.logger.WithField(dlog.KeyOrderID, convertSegment.GetOrderID()).
		WithField(dlog.KeySegmentID, convertSegment.GetID()).
		Info("Publishing segment")

	convertSegment.State = models.SegmentStatePublished

	return d.registry.PersistSegment(convertSegment)
}

// AllocatePublisherAuthority _
func (d *Dealer) AllocatePublisherAuthority(name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/publishers/%s", name)

	return &models.Author{Name: authorName}, nil
}

// GetQueuedSegmentsCount _
func (d *Dealer) GetQueuedSegmentsCount(fctx context.Context) (int, error) {
	segments, err := d.registry.SearchAllSegments(fctx, func(segment models.ISegment) bool {
		return segment.GetState() == models.SegmentStatePublished && !segment.GetIsLocked()
	})

	if err != nil {
		return 0, errors.Wrap(err, "Searching segments")
	}

	return len(segments), nil
}

// GetSegmentsStatesByOrderID _
func (d *Dealer) GetSegmentsStatesByOrderID(fctx context.Context, orderID string) (map[string]string, error) {
	segments, err := d.registry.FindSegmentsByOrderID(fctx, orderID)

	if err != nil {
		return nil, errors.Wrap(err, "Getting segments")
	}

	if len(segments) == 0 {
		return nil, models.ErrNotFound
	}

	statesMap := make(map[string]string)

	for _, segment := range segments {
		statesMap[segment.GetID()] = segment.GetState()
	}

	return statesMap, nil
}

// GetSegmentByID _
func (d *Dealer) GetSegmentByID(segmentID string) (models.ISegment, error) {
	return d.registry.FindSegmentByID(segmentID)
}
