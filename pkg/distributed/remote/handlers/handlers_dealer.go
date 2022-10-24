package handlers

import (
	"context"

	"github.com/pkg/errors"
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

// Notify handles gRPC requests
func (h *DealerHandler) Notify(ctx context.Context, notification *pb.ProgressNotification) (*pb.Empty, error) {
	_, mProgress, err := converters.FromRPCProgress(notification)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileCastingProgressStep(err))
	}

	switch mProgress.Step() {
	case models.UploadingInputStep:
		err = h.dealer.NotifyRawUpload(ctx, models.LocalAuthor, notification.SegmentId, mProgress)

	case models.DownloadingInputStep:
		err = h.dealer.NotifyRawDownload(ctx, models.LocalAuthor, notification.SegmentId, mProgress)

	case models.ProcessingStep:
		err = h.dealer.NotifyProcess(ctx, models.LocalAuthor, notification.SegmentId, mProgress)

	case models.UploadingOutputStep:
		err = h.dealer.NotifyResultUpload(ctx, models.LocalAuthor, notification.SegmentId, mProgress)

	case models.DownloadingOutputStep:
		err = h.dealer.NotifyResultDownload(ctx, models.LocalAuthor, notification.SegmentId, mProgress)

	default:
		err = errors.Wrapf(models.ErrNotImplemented, "Received unimplemented step: `%s`", mProgress.Step())
	}

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// PublishSegment _
func (h *DealerHandler) PublishSegment(ctx context.Context, gReq *pb.PublishSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.PublishSegment(ctx, models.LocalAuthor, gReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// RepublishSegment _
func (h *DealerHandler) RepublishSegment(ctx context.Context, gReq *pb.RepublishSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.RepublishSegment(ctx, models.LocalAuthor, gReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// AcceptSegment _
func (h *DealerHandler) AcceptSegment(ctx context.Context, gReq *pb.AcceptSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.AcceptSegment(ctx, models.LocalAuthor, gReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// FinishSegment _
func (h *DealerHandler) FinishSegment(ctx context.Context, gReq *pb.FinishSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.FinishSegment(ctx, models.LocalAuthor, gReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// QuitSegment _
func (h *DealerHandler) QuitSegment(ctx context.Context, gReq *pb.QuitSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.QuitSegment(ctx, models.LocalAuthor, gReq.SegmentId)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// CancelSegment _
func (h *DealerHandler) CancelSegment(ctx context.Context, gReq *pb.CancelSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.CancelSegment(ctx, models.LocalAuthor, gReq.SegmentId, gReq.Reason)

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}

// FailSegment _
func (h *DealerHandler) FailSegment(ctx context.Context, gReq *pb.FailSegmentRequest) (*pb.Empty, error) {
	err := h.dealer.FailSegment(ctx, models.LocalAuthor, gReq.SegmentId, errors.New(gReq.Failure))

	if err != nil {
		return nil, converters.ToRPCError(errs.WhileHandleRequest(err))
	}

	return &pb.Empty{}, nil
}
