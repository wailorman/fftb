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
	grpcClient    pb.DealerClient
	storageClient models.IStorageClient
}

// NewDealer _
func NewDealer(grpcClient pb.DealerClient, storageClient models.IStorageClient) *Dealer {
	return &Dealer{
		grpcClient:    grpcClient,
		storageClient: storageClient,
	}
}

// AllocateSegment implements IDealer interface
func (d *Dealer) AllocateSegment(ctx context.Context, publisher models.IAuthor, mReq models.IDealerRequest) (models.ISegment, error) {
	gReq, err := converters.ToRPCDealerRequest(models.LocalAuthor.Name, mReq)

	if err != nil {
		return nil, errs.WhileSerializeRequest(err)
	}

	gSegment, err := d.grpcClient.AllocateSegment(ctx, gReq)

	if err != nil {
		return nil, errs.WhilePerformRequest(converters.FromRPCError(err))
	}

	mSegment, err := converters.FromRPCSegment(gSegment)

	if err != nil {
		return nil, errs.WhileDeserializeResponse(err)
	}

	return mSegment, nil
}

// GetOutputStorageClaim implements IDealer interface
func (d *Dealer) GetOutputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	scReq := converters.ToRPCStorageClaimRequest(publisher.GetName(), segmentID)

	rpcStorageClaim, err := d.grpcClient.GetOutputStorageClaim(ctx, scReq)

	if err != nil {
		return nil, errs.WhileGetOutputStorageClaim(converters.FromRPCError(err))
	}

	localStorageClaim, err := d.storageClient.BuildStorageClaimByURL(rpcStorageClaim.Url)

	if err != nil {
		return nil, errs.WhileBuildStorageClaimByURL(err)
	}

	return localStorageClaim, nil
}

// AllocateInputStorageClaim implements IDealer interface
func (d *Dealer) AllocateInputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	scReq := converters.ToRPCStorageClaimRequest(publisher.GetName(), segmentID)

	rpcStorageClaim, err := d.grpcClient.AllocateInputStorageClaim(ctx, scReq)

	if err != nil {
		return nil, errs.WhileAllocateInputStorageClaim(converters.FromRPCError(err))
	}

	localStorageClaim, err := d.storageClient.BuildStorageClaimByURL(rpcStorageClaim.Url)

	if err != nil {
		return nil, errs.WhileBuildStorageClaimByURL(err)
	}

	return localStorageClaim, nil
}

// GetInputStorageClaim implements IDealer interface
func (d *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	scReq := converters.ToRPCStorageClaimRequest(performer.GetName(), segmentID)

	rpcStorageClaim, err := d.grpcClient.GetInputStorageClaim(ctx, scReq)

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

	rpcStorageClaim, err := d.grpcClient.AllocateOutputStorageClaim(ctx, scReq)

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

	_, err = d.grpcClient.Notify(ctx, gProgress)

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

// PublishSegment _
func (d *Dealer) PublishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.grpcClient.PublishSegment(ctx, &pb.PublishSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// RepublishSegment _
func (d *Dealer) RepublishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.grpcClient.RepublishSegment(ctx, &pb.RepublishSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// AcceptSegment _
func (d *Dealer) AcceptSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.grpcClient.AcceptSegment(ctx, &pb.AcceptSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// FinishSegment _
func (d *Dealer) FinishSegment(ctx context.Context, publisher models.IAuthor, segmentID string) error {
	_, err := d.grpcClient.FinishSegment(ctx, &pb.FinishSegmentRequest{
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
	_, err := d.grpcClient.QuitSegment(ctx, &pb.QuitSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// CancelSegment _
func (d *Dealer) CancelSegment(ctx context.Context, publisher models.IAuthor, segmentID string, cancellationReason string) error {
	_, err := d.grpcClient.CancelSegment(ctx, &pb.CancelSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
		Reason:        cancellationReason,
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}

// FailSegment _
func (d *Dealer) FailSegment(ctx context.Context, publisher models.IAuthor, segmentID string, failure error) error {
	_, err := d.grpcClient.FailSegment(ctx, &pb.FailSegmentRequest{
		Authorization: publisher.GetName(),
		SegmentId:     segmentID,
		Failure:       "failure.Error()",
	})

	if err != nil {
		return converters.FromRPCError(errs.WhilePerformRequest(err))
	}

	return nil
}
