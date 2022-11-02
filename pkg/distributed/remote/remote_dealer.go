package remote

import (
	"context"

	"github.com/wailorman/fftb/pkg/distributed/errs"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/converters"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

// Dealer _
type Dealer struct {
	models.IDealer
	rpcClient     pb.Dealer
	storageClient models.IStorageClient
}

// NewDealer _
func NewDealer(rpcClient pb.Dealer, storageClient models.IStorageClient) *Dealer {
	return &Dealer{
		rpcClient:     rpcClient,
		storageClient: storageClient,
	}
}

func (d *Dealer) AllocatePerformerAuthority(ctx context.Context, name string) (models.IAuthor, error) {
	return models.LocalAuthor, nil
}

// GetInputStorageClaim implements IDealer interface
func (d *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	scReq := converters.ToRPCStorageClaimRequest(performer.GetName(), segmentID)

	rpcStorageClaim, err := d.rpcClient.GetInputStorageClaim(ctx, scReq)

	if err != nil {
		return nil, errs.WhileGetInputStorageClaim(converters.FromRPCError(err))
	}

	localStorageClaim, err := d.storageClient.BuildStorageClaimByURL(rpcStorageClaim.Url)

	if err != nil {
		return nil, errs.WhileBuildStorageClaimByURL(err)
	}

	return localStorageClaim, nil
}

// AllocateOutputStorageClaim implements IDealer interface
func (d *Dealer) AllocateOutputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	scReq := converters.ToRPCStorageClaimRequest(performer.GetName(), segmentID)

	rpcStorageClaim, err := d.rpcClient.AllocateOutputStorageClaim(ctx, scReq)

	if err != nil {
		return nil, errs.WhileAllocateOutputStorageClaim(converters.FromRPCError(err))
	}

	localStorageClaim, err := d.storageClient.BuildStorageClaimByURL(rpcStorageClaim.Url)

	if err != nil {
		return nil, errs.WhileBuildStorageClaimByURL(err)
	}

	return localStorageClaim, nil
}

func (d *Dealer) notify(ctx context.Context, publisher models.IAuthor, segmentID string, p models.IProgress) error {
	gProgress, err := converters.ToRPCProgress(models.LocalAuthor.Name, segmentID, p)

	if err != nil {
		return errs.WhileSerializeRequest(err)
	}

	_, err = d.rpcClient.Notify(ctx, gProgress)

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// NotifyRawUpload implements IDealer interface
func (d *Dealer) NotifyRawUpload(ctx context.Context, publisher models.IAuthor, segmentID string, p models.IProgress) error {
	return d.notify(ctx, publisher, segmentID, p)
}

// NotifyRawDownload implements IDealer interface
func (d *Dealer) NotifyRawDownload(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.notify(ctx, performer, segmentID, p)
}

// NotifyProcess implements IDealer interface
func (d *Dealer) NotifyProcess(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.notify(ctx, performer, segmentID, p)
}

// NotifyResultUpload implements IDealer interface
func (d *Dealer) NotifyResultUpload(ctx context.Context, performer models.IAuthor, segmentID string, p models.IProgress) error {
	return d.notify(ctx, performer, segmentID, p)
}

// NotifyResultDownload implements IDealer interface
func (d *Dealer) NotifyResultDownload(ctx context.Context, publisher models.IAuthor, segmentID string, p models.IProgress) error {
	return d.notify(ctx, publisher, segmentID, p)
}

// FinishSegment _
func (d *Dealer) FinishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.rpcClient.FinishSegment(ctx, &pb.FinishSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// QuitSegment _
func (d *Dealer) QuitSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.rpcClient.QuitSegment(ctx, &pb.QuitSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// FailSegment _
func (d *Dealer) FailSegment(ctx context.Context, publisher models.IAuthor, segmentID string, failure error) error {
	_, err := d.rpcClient.FailSegment(ctx, &pb.FailSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
		Failure:       failure.Error(),
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

func (d *Dealer) FindFreeSegment(ctx context.Context, performer models.IAuthor) (models.ISegment, error) {
	rpcSegment, err := d.rpcClient.FindFreeSegment(ctx, &pb.FindFreeSegmentRequest{Authorization: performer.GetName()})

	if err != nil {
		return nil, converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	mSeg, err := converters.FromRPCSegment(rpcSegment)

	if err != nil {
		return nil, errs.WhileDeserializeResponse(err)
	}

	return mSeg, nil
}
