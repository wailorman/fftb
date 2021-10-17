package handlers

import (
	"context"

	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

// DealerHandler _
type DealerHandler struct {
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

// AllocateSegment _
func (h *DealerHandler) AllocateSegment(ctx context.Context, req *pb.AllocateSegmentRequest) (*pb.AllocateSegmentResponse, error) {
	h.dealer.AllocateSegment(ctx)
}
