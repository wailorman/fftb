package handlers

import (
	"context"

	"github.com/wailorman/fftb/pkg/distributed/errs"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/converters"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

// DealerHandler _
type DealerHandler struct {
	pb.DealerServer
	ctx             context.Context
	dealer          models.IDealer
	contracter      models.IContracter
	authoritySecret []byte
	sessionSecret   []byte
}

// NewDealerHandler _
func NewDealerHandler(
	dealer models.IDealer,
	authoritySecret []byte,
	sessionSecret []byte) *DealerHandler {

	// TODO: handler config

	return &DealerHandler{
		dealer:          dealer,
		authoritySecret: authoritySecret,
		sessionSecret:   sessionSecret,
	}
}

// AllocateSegment handles gRPC requests
func (h *DealerHandler) AllocateSegment(ctx context.Context, gReq *pb.DealerRequest) (*pb.Segment, error) {
	_, mReq, err := converters.FromRPCDealerRequest(gReq)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileDeserializeRequest(err))
	}

	mSeg, err := h.dealer.AllocateSegment(ctx, models.LocalAuthor, mReq)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	gSeg, err := converters.ToRPCSegment(mSeg)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileSerializeResponse(err))
	}

	return gSeg, nil
}

// GetOutputStorageClaim handles gRPC requests
func (h *DealerHandler) GetOutputStorageClaim(ctx context.Context, scReq *pb.StorageClaimRequest) (*pb.StorageClaim, error) {
	storageClaim, err := h.dealer.GetOutputStorageClaim(ctx, models.LocalAuthor, scReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileGetOutputStorageClaim(err))
	}

	return converters.ToRPCStorageClaim(storageClaim), nil
}

// AllocateInputStorageClaim handles gRPC requests
func (h *DealerHandler) AllocateInputStorageClaim(ctx context.Context, scReq *pb.StorageClaimRequest) (*pb.StorageClaim, error) {
	storageClaim, err := h.dealer.AllocateInputStorageClaim(ctx, models.LocalAuthor, scReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileAllocateInputStorageClaim(err))
	}

	return converters.ToRPCStorageClaim(storageClaim), nil
}

// GetInputStorageClaim handles gRPC requests
func (h *DealerHandler) GetInputStorageClaim(ctx context.Context, scReq *pb.StorageClaimRequest) (*pb.StorageClaim, error) {
	storageClaim, err := h.dealer.GetInputStorageClaim(ctx, models.LocalAuthor, scReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileGetInputStorageClaim(err))
	}

	return converters.ToRPCStorageClaim(storageClaim), nil
}

// AllocateOutputStorageClaim handles gRPC requests
func (h *DealerHandler) AllocateOutputStorageClaim(ctx context.Context, scReq *pb.StorageClaimRequest) (*pb.StorageClaim, error) {
	storageClaim, err := h.dealer.AllocateOutputStorageClaim(ctx, models.LocalAuthor, scReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileAllocateOutputStorageClaim(err))
	}

	return converters.ToRPCStorageClaim(storageClaim), nil
}
