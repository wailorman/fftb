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

// AllocateSegment implements IDealer method
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

// GetOutputStorageClaim implements IDealer method
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

// AllocateInputStorageClaim implements IDealer method
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

// GetInputStorageClaim implements IDealer method
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

// AllocateOutputStorageClaim implements IDealer method
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
