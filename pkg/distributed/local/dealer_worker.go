package local

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// FindFreeSegment _
func (d *Dealer) FindFreeSegment(ctx context.Context, performer models.IAuthor) (models.ISegment, error) {
	if !d.freeSegmentLock.TryLockTimeout(LockSegmentTimeout) {
		return nil, models.ErrLockTimeout
	}

	defer d.freeSegmentLock.Unlock()

	segment, err := d.registry.SearchSegment(ctx, func(segment models.ISegment) bool {
		return segment.GetCanPerform()
	})

	if err != nil {
		return nil, errors.Wrap(err, "Looking for free segment")
	}

	err = d.segmentMutator.LockSegment(segment, performer)

	if err != nil {
		return nil, errors.Wrap(err, "Locking segment")
	}

	err = d.registry.PersistSegment(ctx, segment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting segment")
	}

	return segment, nil
}

// GetInputStorageClaim _
func (d *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	// TODO: match performer
	return d.getInputStorageClaim(ctx, segmentID)
}

// AllocateOutputStorageClaim _
func (d *Dealer) AllocateOutputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownType
	}

	oClaimID := fmt.Sprintf("output_%s_%s_%s", segment.GetOrderID(), segment.GetID(), uuid.New().String())
	oClaim, err := d.storageController.AllocateStorageClaim(ctx, oClaimID)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating output storage claim")
	}

	convertSegment.OutputStorageClaimIdentity = oClaimID

	err = d.registry.PersistSegment(ctx, convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting output claim identity")
	}

	return oClaim, nil
}

// FinishSegment _
func (d *Dealer) FinishSegment(ctx context.Context, performer models.IAuthor, segmentID string) error {
	segment, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownType
	}

	dlog.WithSegment(d.logger, segment).
		Info("Segment is finished")

	err = d.segmentMutator.FinishSegment(convertSegment)

	if err != nil {
		return errors.Wrap(err, "Finishing segment")
	}

	err = d.registry.PersistSegment(ctx, convertSegment)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}

// NotifyProcess _
func (d *Dealer) NotifyProcess(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.segmentProgress(ctx, performer, segmentID, p)
}

// NotifyRawDownload _
func (d *Dealer) NotifyRawDownload(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.segmentProgress(ctx, performer, segmentID, p)
}

// NotifyResultUpload _
func (d *Dealer) NotifyResultUpload(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.segmentProgress(ctx, performer, segmentID, p)
}

func (d *Dealer) segmentProgress(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	seg, err := d.registry.FindSegmentByID(ctx, segmentID)

	if err != nil {
		return err
	}

	dlog.SegmentProgress(d.logger, seg, p)

	err = d.segmentMutator.LockSegment(seg, performer)

	if err != nil {
		return errors.Wrap(err, "Locking segment")
	}

	err = d.registry.PersistSegment(ctx, seg)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}

// AllocatePerformerAuthority _
func (d *Dealer) AllocatePerformerAuthority(ctx context.Context, name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/performers/%s", name)

	return &models.Author{Name: authorName}, nil
}

// // WaitOnSegmentFinished _
// func (d *Dealer) WaitOnSegmentFinished(ctx context.Context, id string) <-chan struct{} {
// 	panic("not implemented") // TODO:
// }

// // WaitOnSegmentFailed _
// func (d *Dealer) WaitOnSegmentFailed(ctx context.Context, id string) <-chan error {
// 	panic("not implemented") // TODO:
// }

// // WaitOnSegmentCancelled _
// func (d *Dealer) WaitOnSegmentCancelled(ctx context.Context, id string) <-chan struct{} {
// 	panic("not implemented") // TODO:
// }

// FailSegment _
func (d *Dealer) FailSegment(ctx context.Context, performer models.IAuthor, id string, reportedErr error) error {
	segment, err := d.registry.FindSegmentByID(ctx, id)

	if err != nil {
		return errors.Wrapf(err, "Searching segment by id `%s`", id)
	}

	dlog.WithSegment(d.logger, segment).
		WithError(reportedErr).
		Info("Received segment failure")

	err = d.segmentMutator.FailSegment(segment, reportedErr)

	if err != nil {
		return errors.Wrap(err, "Failing segment")
	}

	err = d.registry.PersistSegment(ctx, segment)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}

// QuitSegment _
func (d *Dealer) QuitSegment(ctx context.Context, performer models.IAuthor, id string) error {
	seg, err := d.registry.FindSegmentByID(ctx, id)

	if err != nil {
		return err
	}

	dlog.WithSegment(d.logger, seg).
		WithField(dlog.KeyPerformer, performer.GetName()).
		Debug("Quitting segment")

	err = d.segmentMutator.UnlockSegment(seg)

	if err != nil {
		return errors.Wrap(err, "Unlocking segment")
	}

	err = d.registry.PersistSegment(ctx, seg)

	if err != nil {
		return errors.Wrap(err, "Persisting segment")
	}

	return nil
}
