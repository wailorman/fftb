package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// DealerHandler _
type DealerHandler struct {
	ctx    context.Context
	dealer models.IDealer
}

// NewDealerHandler _
func NewDealerHandler(ctx context.Context, dealer models.IDealer) *DealerHandler {
	return &DealerHandler{
		ctx:    ctx,
		dealer: dealer,
	}
}

// AllocatePublisherAuthority _
// // POST /publisher_authorities
func (dh *DealerHandler) AllocatePublisherAuthority(c *gin.Context) {

}

// AllocatePerformerAuthority _
// // POST /performer_authorities
func (dh *DealerHandler) AllocatePerformerAuthority(c *gin.Context) {
}

// FindFreeSegment _
// // POST /segments/free | Segment
func (dh *DealerHandler) FindFreeSegment(c *gin.Context) {
}

// GetSegmentByID _
// // GET /segments/{id} | Segment
func (dh *DealerHandler) GetSegmentByID(c *gin.Context) {
}

// AllocateSegment _
// // POST /segments
func (dh *DealerHandler) AllocateSegment(c *gin.Context) {
}

// PublishSegment _
// // POST /segments/{id}/actions/publish
func (dh *DealerHandler) PublishSegment(c *gin.Context) {
}

// RepublishSegment _
// // POST /segments/{id}/actions/republish
func (dh *DealerHandler) RepublishSegment(c *gin.Context) {
}

// CancelSegment _
// // POST /segments/{id}/actions/cancel | { reason: failed }
func (dh *DealerHandler) CancelSegment(c *gin.Context) {
}

// AcceptSegment _
// // POST /segments/{id}/actions/accept
func (dh *DealerHandler) AcceptSegment(c *gin.Context) {
}

// FinishSegment _
// // POST /segments/{id}/actions/finish
func (dh *DealerHandler) FinishSegment(c *gin.Context) {
}

// QuitSegment _
// // POST /segments/{id}/actions/quit
func (dh *DealerHandler) QuitSegment(c *gin.Context) {
}

// FailSegment _
// // POST /segments/{id}/actions/fail
func (dh *DealerHandler) FailSegment(c *gin.Context) {
}

// GetOutputStorageClaim _
// // GET /segments/{id}/output_storage_claim | { storage_claim: http... }
func (dh *DealerHandler) GetOutputStorageClaim(c *gin.Context) {
}

// GetInputStorageClaim _
// // GET /segments/{id}/input_storage_claim | { storage_claim: http... }
func (dh *DealerHandler) GetInputStorageClaim(c *gin.Context) {
}

// AllocateInputStorageClaim _
// // POST /segments/{id}/output_storage_claim | { storage_claim: http... }
func (dh *DealerHandler) AllocateInputStorageClaim(c *gin.Context) {
}

// AllocateOutputStorageClaim _
// // POST /segments/{id}/input_storage_claim | { storage_claim: http... }
func (dh *DealerHandler) AllocateOutputStorageClaim(c *gin.Context) {
}

// NotifyRawUpload _
// // POST /segments/{id}/notifications/input_upload | { progress: 0.5 }
func (dh *DealerHandler) NotifyRawUpload(c *gin.Context) {
}

// NotifyResultDownload _
// // POST /segments/{id}/notifications/output_download | { progress: 0.5 }
func (dh *DealerHandler) NotifyResultDownload(c *gin.Context) {
}

// NotifyRawDownload _
// // POST /segments/{id}/notifications/input_download | { progress: 0.5 }
func (dh *DealerHandler) NotifyRawDownload(c *gin.Context) {
}

// NotifyResultUpload _
// // POST /segments/{id}/notifications/ouput_upload | { progress: 0.5 }
func (dh *DealerHandler) NotifyResultUpload(c *gin.Context) {
}

// NotifyProcess _
// // POST /segments/{id}/notifications/process | { progress: 0.5 }
func (dh *DealerHandler) NotifyProcess(c *gin.Context) {
}
