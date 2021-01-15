package local

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// Dealer _
type Dealer struct {
	storageController models.IStorageController
	registry          models.IRegistry
	freeSegmentLock   trylock.TryLocker
	logger            logrus.FieldLogger
	ctx               context.Context
}

// NewDealer _
func NewDealer(ctx context.Context, sc models.IStorageController, r models.IRegistry) *Dealer {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.distributed.dealer"); logger == nil {
		logger = ctxlog.New("fftb.distributed.dealer")
	}

	return &Dealer{
		storageController: sc,
		registry:          r,
		freeSegmentLock:   trylock.New(),
		logger:            logger,
		ctx:               ctx,
	}
}

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

	return convertSegment, nil
}

// FindFreeSegment _
func (d *Dealer) FindFreeSegment() (models.ISegment, error) {
	locked := d.freeSegmentLock.TryLockTimeout(LockSegmentTimeout)

	if !locked {
		return nil, models.ErrFreeSegmentLockTimeout
	}

	defer d.freeSegmentLock.Unlock()

	freeSegment, err := d.registry.FindNotLockedSegment()

	if err != nil {
		return nil, errors.Wrap(err, "Looking for free segment")
	}

	err = d.registry.LockSegmentByID(freeSegment.GetID(), performer)

	if err != nil {
		return nil, errors.Wrap(err, "Locking free segment")
	}

	return freeSegment, nil
}

// GetInputStorageClaim _
func (d *Dealer) GetInputStorageClaim(performer IAuthor, segment models.ISegment) (models.IStorageClaim, error) {

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

// GetOutputStorageClaim _
func (d *Dealer) GetOutputStorageClaim(segment models.ISegment) (models.IStorageClaim, error) {
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
func (d *Dealer) AllocateInputStorageClaim(publisher models.IAuthor, segment models.ISegment) (models.IStorageClaim, error) {
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

// AllocateOutputStorageClaim _
func (d *Dealer) AllocateOutputStorageClaim(performer models.IAuthor, segment models.ISegment) (models.IStorageClaim, error) {
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

// CancelSegment _
func (d *Dealer) CancelSegment(publisher models.IAuthor, seg models.ISegment) error {
	panic(models.ErrNotImplemented)
}

// FindSegmentByID _
func (d *Dealer) FindSegmentByID(id string) (models.ISegment, error) {
	panic(models.ErrNotImplemented)
}

// NotifyRawUpload _
func (d *Dealer) NotifyRawUpload(publisher models.IAuthor, seg models.ISegment, p models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// NotifyResultDownload _
func (d *Dealer) NotifyResultDownload(publisher models.IAuthor, seg models.ISegment, p models.Progresser) error {
	panic(models.ErrNotImplemented)
}

// PublishSegment _
func (d *Dealer) PublishSegment(publisher models.IAuthor, segment models.ISegment) error {
	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	convertSegment.State = models.SegmentPublishedState

	return d.registry.PersistSegment(convertSegment)
}

// // Subscription _
// func (d *Dealer) Subscription(segment models.ISegment) (models.Subscriber, error) {
// 	panic(models.ErrNotImplemented)
// }

// FinishSegment _
func (d *Dealer) FinishSegment(performer models.IAuthor, segment models.ISegment) error {
	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	convertSegment.State = models.SegmentFinishedState

	return d.registry.PersistSegment(convertSegment)
}

// NotifyProcess _
func (d *Dealer) NotifyProcess(performer models.IAuthor, seg models.ISegment, p models.Progresser) error {
	return d.segmentProgress(performer, seg, p)
}

// NotifyRawDownload _
func (d *Dealer) NotifyRawDownload(performer models.IAuthor, seg models.ISegment, p models.Progresser) error {
	return d.segmentProgress(performer, seg, p)
}

// NotifyResultUpload _
func (d *Dealer) NotifyResultUpload(performer models.IAuthor, seg models.ISegment, p models.Progresser) error {
	return d.segmentProgress(performer, seg, p)
}

func (d *Dealer) segmentProgress(performer models.IAuthor, seg models.ISegment, p models.Progresser) error {
	dlog.SegmentProgress(d.logger, seg, p)

	err := d.registry.LockSegmentByID(seg.GetID(), performer)

	if err != nil {
		return errors.Wrap(err, "Prolongating segment lock")
	}

	return nil
}

// AllocatePublisherAuthority _
func (d *Dealer) AllocatePublisherAuthority(name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/publishers/%s", name)

	return &models.Author{Name: authorName}, nil
}

// AllocatePerformerAuthority _
func (d *Dealer) AllocatePerformerAuthority(name string) (models.IAuthor, error) {
	authorName := fmt.Sprintf("v1/performers/%s", name)

	return &models.Author{Name: authorName}, nil
}
