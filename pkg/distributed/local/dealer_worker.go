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
func (d *Dealer) FindFreeSegment(performer models.IAuthor) (models.ISegment, error) {
	if !d.freeSegmentLock.TryLockTimeout(LockSegmentTimeout) {
		return nil, models.ErrFreeSegmentLockTimeout
	}

	defer d.freeSegmentLock.Unlock()

	freeSegment, err := d.registry.FindNotLockedSegment(d.ctx)

	if err != nil {
		return nil, errors.Wrap(err, "Looking for free segment")
	}

	freeSegment.Lock(performer)

	err = d.registry.PersistSegment(freeSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting locked segment")
	}

	return freeSegment, nil
}

// GetInputStorageClaim _
func (d *Dealer) GetInputStorageClaim(performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	if convertSegment.InputStorageClaimIdentity == "" {
		return nil, errors.Wrap(models.ErrMissingStorageClaim, "Getting input storage claim identity")
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.InputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

// AllocateOutputStorageClaim _
func (d *Dealer) AllocateOutputStorageClaim(performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	oClaimID := fmt.Sprintf("%s/%s/output_%s", segment.GetOrderID(), segment.GetID(), uuid.New().String())
	oClaim, err := d.storageController.AllocateStorageClaim(oClaimID)

	if err != nil {
		return nil, errors.Wrap(err, "Allocating output storage claim")
	}

	convertSegment.OutputStorageClaimIdentity = oClaimID

	err = d.registry.PersistSegment(convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting output claim identity")
	}

	return oClaim, nil
}

// FinishSegment _
func (d *Dealer) FinishSegment(performer models.IAuthor, segmentID string) error {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	d.logger.WithField(dlog.KeySegmentID, convertSegment.GetID()).
		WithField(dlog.KeyOrderID, segment.GetOrderID()).
		Info("Segment is finished")

	convertSegment.State = models.SegmentStateFinished
	convertSegment.Unlock()

	return d.registry.PersistSegment(convertSegment)
}

// NotifyProcess _
func (d *Dealer) NotifyProcess(performer models.IAuthor, segmentID string, p models.Progresser) error {
	return d.segmentProgress(performer, segmentID, p)
}

// NotifyRawDownload _
func (d *Dealer) NotifyRawDownload(performer models.IAuthor, segmentID string, p models.Progresser) error {
	return d.segmentProgress(performer, segmentID, p)
}

// NotifyResultUpload _
func (d *Dealer) NotifyResultUpload(performer models.IAuthor, segmentID string, p models.Progresser) error {
	return d.segmentProgress(performer, segmentID, p)
}

func (d *Dealer) segmentProgress(performer models.IAuthor, segmentID string, p models.Progresser) error {
	seg, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return err
	}

	dlog.SegmentProgress(d.logger, seg, p)

	seg.Lock(performer)

	err = d.registry.PersistSegment(seg)

	if err != nil {
		return errors.Wrap(err, "Prolongating segment lock")
	}

	return nil
}

// AllocatePerformerAuthority _
func (d *Dealer) AllocatePerformerAuthority(name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/performers/%s", name)

	return &models.Author{Name: authorName}, nil
}

// WaitOnSegmentFinished _
func (d *Dealer) WaitOnSegmentFinished(ctx context.Context, id string) <-chan struct{} {
	panic("not implemented")
}

// WaitOnSegmentFailed _
func (d *Dealer) WaitOnSegmentFailed(ctx context.Context, id string) <-chan error {
	panic("not implemented")
}

// WaitOnSegmentCancelled _
func (d *Dealer) WaitOnSegmentCancelled(ctx context.Context, id string) <-chan struct{} {
	panic("not implemented")
}

// FailSegment _
func (d *Dealer) FailSegment(performer models.IAuthor, id string, err error) error {
	panic("not implemented")
}

// QuitSegment _
func (d *Dealer) QuitSegment(performer models.IAuthor, id string) error {
	d.logger.WithField(dlog.KeyPerformer, performer.GetName()).
		WithField(dlog.KeySegmentID, id).
		Debug("Quitting segment")

	seg, err := d.registry.FindSegmentByID(id)

	if err != nil {
		return err
	}

	if seg.GetPerformer() == nil {
		return nil
	}

	if !seg.GetPerformer().IsEqual(performer) {
		return errors.Wrap(models.ErrPerformerMismatch, fmt.Sprintf("Received performer `%s`, locked by performer `%s`", performer, seg.GetPerformer()))
	}

	seg.Unlock()

	return d.registry.PersistSegment(seg)
}
