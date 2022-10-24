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
func (d *Dealer) AllocateSegment(ctx context.Context, publisher models.IAuthor, req models.IDealerRequest) (models.ISegment, error) {
	convertReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return nil, models.ErrUnknownType
	}

	if validationErr := req.Validate(); validationErr != nil {
		return nil, errors.Wrap(validationErr, "Validation error")
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
		Publisher:     publisher,
	}

	err := d.registry.PersistSegment(ctx, convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting segment")
	}

	return convertSegment, nil
}

// GetOutputStorageClaim _
func (d *Dealer) GetOutputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	// TODO: match publisher
	return d.getOutputStorageClaim(ctx, segmentID)
}

// AllocateInputStorageClaim _
func (d *Dealer) AllocateInputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return nil, errors.Wrapf(err, "Finding segment by id `%s`", segmentID)
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownType
	}

	iClaimID := fmt.Sprintf("input_%s_%s_%s", segment.GetOrderID(), segment.GetID(), uuid.New().String())
	iClaim, err := d.storageController.AllocateStorageClaim(ctx, iClaimID)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating input storage claim")
	}

	convertSegment.InputStorageClaimIdentity = iClaimID

	err = d.registry.PersistSegment(ctx, convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting input claim identity")
	}

	return iClaim, nil
}

// CancelSegment _
func (d *Dealer) CancelSegment(ctx context.Context, publisher models.IAuthor, segmentID string, reason string) error {
	// TODO: lock segment
	// TODO: receive multiple segment ids
	// TODO: match publisher

	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return errors.Wrapf(err, "Finding segment by id `%s`", segmentID)
	}

	// convertSegment, ok := segment.(*models.ConvertSegment)

	// if !ok {
	// 	return models.ErrUnknownType
	// }

	dlog.WithSegment(d.logger, segment).
		WithField(dlog.KeyReason, reason).
		Info("Cancelling segment")

	if segment.GetState() == models.SegmentStateCancelled {
		return nil
	}

	err = d.segmentMutator.CancelSegment(segment, reason)

	if err != nil {
		return errors.Wrap(err, "Cancelling segment")
	}

	return d.registry.PersistSegment(ctx, segment)
}

// AcceptSegment _
func (d *Dealer) AcceptSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	// TODO: lock segment

	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return errors.Wrapf(err, "Finding segment by id `%s`", segmentID)
	}

	logger := dlog.WithSegment(d.logger, segment)

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownType
	}

	logger.Info("Accepting segment")

	convertSegment.State = models.SegmentStateAccepted

	err = d.registry.PersistSegment(ctx, convertSegment)

	if err != nil {
		return errors.Wrapf(err, "Persisting segment `%s`", segmentID)
	}

	d.tryPurgeInputStorageClaim(segmentID)
	d.tryPurgeOutputStorageClaim(segmentID)

	return nil
}

// NotifyRawUpload _
func (d *Dealer) NotifyRawUpload(ctx context.Context, publisher models.IAuthor, segmentID string, p models.IProgress) error {
	panic(models.ErrNotImplemented)
}

// NotifyResultDownload _
func (d *Dealer) NotifyResultDownload(ctx context.Context, publisher models.IAuthor, segmentID string, p models.IProgress) error {
	panic(models.ErrNotImplemented)
}

// PublishSegment _
func (d *Dealer) PublishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	// TODO: lock segment

	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return errors.Wrapf(err, "Finding segment by id `%s`", segmentID)
	}

	dlog.WithSegment(d.logger, segment).
		Info("Publishing segment")

	err = d.segmentMutator.PublishSegment(segment)

	if err != nil {
		return errors.Wrap(err, "Publishing segment")
	}

	err = d.registry.PersistSegment(ctx, segment)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}

// RepublishSegment _
func (d *Dealer) RepublishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	// TODO: lock segment

	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return errors.Wrapf(err, "Finding segment by id `%s`", segmentID)
	}

	dlog.WithSegment(d.logger, segment).
		Info("Republishing segment")

	err = d.segmentMutator.PublishSegment(segment)

	if err != nil {
		return errors.Wrap(err, "Republishing segment")
	}

	err = d.registry.PersistSegment(ctx, segment)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}

// AllocatePublisherAuthority _
func (d *Dealer) AllocatePublisherAuthority(ctx context.Context, name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/publishers/%s", name)

	return &models.Author{Name: authorName}, nil
}

// GetQueuedSegmentsCount _
func (d *Dealer) GetQueuedSegmentsCount(fctx context.Context, publisher models.IAuthor) (int, error) {
	segments, err := d.registry.SearchAllSegments(fctx, func(segment models.ISegment) bool {
		return segment.GetState() == models.SegmentStatePublished &&
			!segment.GetIsLocked() &&
			segment.GetPublisher() != nil &&
			segment.GetPublisher().IsEqual(publisher)
	})

	if err != nil {
		return 0, errors.Wrap(err, "Searching segments")
	}

	return len(segments), nil
}

// GetSegmentsByOrderID _
func (d *Dealer) GetSegmentsByOrderID(ctx context.Context, publisher models.IAuthor, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	segments, err := d.registry.FindSegmentsByOrderID(ctx, orderID)

	if err != nil {
		return nil, err
	}

	resultSegments := make([]models.ISegment, 0)

	for _, segment := range segments {
		if search.Select(segment) {
			resultSegments = append(resultSegments, segment)
		}
	}

	return resultSegments, nil
}

// // GetSegmentsStatesByOrderID _
// func (d *Dealer) GetSegmentsStatesByOrderID(fctx context.Context, orderID string) (map[string]string, error) {
// 	segments, err := d.GetSegmentsByOrderID(fctx, orderID, models.EmptySegmentFilters())

// 	if err != nil {
// 		return nil, errors.Wrap(err, "Getting segments")
// 	}

// 	if len(segments) == 0 {
// 		return nil, models.ErrNotFound
// 	}

// 	statesMap := make(map[string]string)

// 	for _, segment := range segments {
// 		statesMap[segment.GetID()] = segment.GetState()
// 	}

// 	return statesMap, nil
// }

// GetSegmentByID _
func (d *Dealer) GetSegmentByID(ctx context.Context, publisher models.IAuthor, segmentID string) (models.ISegment, error) {
	// TODO: match publisher

	return d.registry.FindSegmentByID(ctx, segmentID)
}
