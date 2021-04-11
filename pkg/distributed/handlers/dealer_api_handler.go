package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"

	// "github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/distributed/remote"
)

// TODO: remove
var localAuthor models.IAuthor = &models.Author{Name: "local"}

// DealerHandler _
type DealerHandler struct {
	ctx    context.Context
	dealer models.IDealer
}

// APIError _
type APIError struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func newAPIError(err error) (code int, body *remote.ProblemDetails) {
	cause := errors.Cause(err)

	code = 422

	if errors.Is(err, models.ErrNotFound) {
		code = 404
	}

	if errors.Is(err, models.ErrUnknownSegmentType) {
		code = 422
	}

	detail := err.Error()
	errType := "github.com/wailorman/fftb"

	return code, &remote.ProblemDetails{
		Type:   &errType,
		Title:  cause.Error(),
		Detail: &detail,
	}
}

func buildConvertSegment(convSeg *models.ConvertSegment) *remote.Segment {
	return &remote.Segment{
		Id: convSeg.Identity,
		// OrderID:  convSeg.OrderIdentity,
		// Type:     convSeg.Type,
		// State:    convSeg.State,
		// Params:   convSeg.Params,
		// Muxer:    convSeg.Muxer,
		// Position: convSeg.Position,
	}
}

// NewDealerHandler _
func NewDealerHandler(ctx context.Context, dealer models.IDealer) *DealerHandler {
	return &DealerHandler{
		ctx:    ctx,
		dealer: dealer,
	}
}

// // AllocateAuthority _
// // // POST /authorities
// func (dh *DealerHandler) AllocateAuthority(c *gin.Context) {
// 	c.JSON(200, &models.RemoteAuthority{
// 		Authority: "local",
// 	})
// }

// // FindFreeSegment _
// // // POST /segments/free | Segment
// func (dh *DealerHandler) FindFreeSegment(c *gin.Context) {
// 	seg, err := dh.dealer.FindFreeSegment(dh.ctx, localAuthor)

// 	if err != nil {
// 		c.JSON(newAPIError(err))
// 		return
// 	}

// 	convSeg, ok := seg.(*models.ConvertSegment)

// 	if !ok {
// 		c.JSON(newAPIError(models.ErrUnknownSegmentType))
// 		return
// 	}

// 	c.JSON(200, buildConvertSegment(convSeg))
// }

// // GetSegmentByID _
// // // GET /segments/{id} | Segment
// func (dh *DealerHandler) GetSegmentByID(c *gin.Context) {
// 	seg, err := dh.dealer.GetSegmentByID(dh.ctx, localAuthor, c.Param("id"))

// 	if err != nil {
// 		c.JSON(newAPIError(err))
// 		return
// 	}

// 	convSeg, ok := seg.(*models.ConvertSegment)

// 	if !ok {
// 		c.JSON(newAPIError(models.ErrUnknownSegmentType))
// 		return
// 	}

// 	c.JSON(200, buildConvertSegment(convSeg))
// }

// AllocateSegment _
// // POST /segments
func (dh *DealerHandler) AllocateSegment(c echo.Context) error {
	params := &models.ConvertDealerRequest{}

	if err := c.Bind(&params); err != nil {
		c.JSON(newAPIError(err))
		return nil
	}

	seg, err := dh.dealer.AllocateSegment(dh.ctx, localAuthor, params)

	if err != nil {
		c.JSON(newAPIError(err))
		return nil
	}

	convSeg, ok := seg.(*models.ConvertSegment)

	if !ok {
		c.JSON(newAPIError(errors.Wrapf(models.ErrUnknownSegmentType, "Received `%s`", seg.GetType())))
		return nil
	}

	c.JSON(200, buildConvertSegment(convSeg))
	return nil
}

// // PublishSegment _
// // // POST /segments/{id}/actions/publish
// func (dh *DealerHandler) PublishSegment(c *gin.Context) {
// }

// // RepublishSegment _
// // // POST /segments/{id}/actions/republish
// func (dh *DealerHandler) RepublishSegment(c *gin.Context) {
// }

// // CancelSegment _
// // // POST /segments/{id}/actions/cancel | { reason: failed }
// func (dh *DealerHandler) CancelSegment(c *gin.Context) {
// }

// // AcceptSegment _
// // // POST /segments/{id}/actions/accept
// func (dh *DealerHandler) AcceptSegment(c *gin.Context) {
// }

// // FinishSegment _
// // // POST /segments/{id}/actions/finish
// func (dh *DealerHandler) FinishSegment(c *gin.Context) {
// }

// // QuitSegment _
// // // POST /segments/{id}/actions/quit
// func (dh *DealerHandler) QuitSegment(c *gin.Context) {
// }

// // FailSegment _
// // // POST /segments/{id}/actions/fail
// func (dh *DealerHandler) FailSegment(c *gin.Context) {
// }

// // GetOutputStorageClaim _
// // // GET /segments/{id}/output_storage_claim | { storage_claim: http... }
// func (dh *DealerHandler) GetOutputStorageClaim(c *gin.Context) {
// }

// // GetInputStorageClaim _
// // // GET /segments/{id}/input_storage_claim | { storage_claim: http... }
// func (dh *DealerHandler) GetInputStorageClaim(c *gin.Context) {
// }

// // AllocateInputStorageClaim _
// // // POST /segments/{id}/output_storage_claim | { storage_claim: http... }
// func (dh *DealerHandler) AllocateInputStorageClaim(c *gin.Context) {
// }

// // AllocateOutputStorageClaim _
// // // POST /segments/{id}/input_storage_claim | { storage_claim: http... }
// func (dh *DealerHandler) AllocateOutputStorageClaim(c *gin.Context) {
// }

// // NotifyRawUpload _
// // // POST /segments/{id}/notifications/input_upload | { progress: 0.5 }
// func (dh *DealerHandler) NotifyRawUpload(c *gin.Context) {
// }

// // NotifyResultDownload _
// // // POST /segments/{id}/notifications/output_download | { progress: 0.5 }
// func (dh *DealerHandler) NotifyResultDownload(c *gin.Context) {
// }

// // NotifyRawDownload _
// // // POST /segments/{id}/notifications/input_download | { progress: 0.5 }
// func (dh *DealerHandler) NotifyRawDownload(c *gin.Context) {
// }

// // NotifyResultUpload _
// // // POST /segments/{id}/notifications/ouput_upload | { progress: 0.5 }
// func (dh *DealerHandler) NotifyResultUpload(c *gin.Context) {
// }

// // NotifyProcess _
// // // POST /segments/{id}/notifications/process | { progress: 0.5 }
// func (dh *DealerHandler) NotifyProcess(c *gin.Context) {
// }
