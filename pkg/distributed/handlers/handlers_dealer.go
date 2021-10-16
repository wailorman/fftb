package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
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

func buildConvertSegment(convSeg *models.ConvertSegment) *dSchema.ConvertSegment {
	return &dSchema.ConvertSegment{
		Type:     models.ConvertV1Type,
		Id:       convSeg.Identity,
		OrderId:  convSeg.OrderIdentity,
		Muxer:    convSeg.Muxer,
		Position: convSeg.Position,
		Params: dSchema.ConvertParams{
			HwAccel:          convSeg.Params.HWAccel,
			KeyframeInterval: convSeg.Params.KeyframeInterval,
			Preset:           convSeg.Params.Preset,
			Scale:            convSeg.Params.Scale,
			VideoBitRate:     convSeg.Params.VideoBitRate,
			VideoCodec:       convSeg.Params.VideoCodec,
			VideoQuality:     convSeg.Params.VideoQuality,
		},
	}
}

// AllocateAuthority _
// POST /authorities
func (dh *DealerHandler) AllocateAuthority(c echo.Context) error {
	params := &dSchema.AuthorityInput{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	key, err := CreateAuthorityToken(dh.authoritySecret, params.Name)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.Authority{Key: key})
}

// CreateSession _
// POST /sessions
func (dh *DealerHandler) CreateSession(c echo.Context) error {
	params := &dSchema.SessionInput{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	key, err := CreateSessionToken(dh.authoritySecret, dh.sessionSecret, params.AuthorityKey)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.Session{Key: key})
}

// FindFreeSegment _
// // POST /segments/free | Segment
func (dh *DealerHandler) FindFreeSegment(c echo.Context) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	seg, err := dh.dealer.FindFreeSegment(c.Request().Context(), author)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	convSeg, ok := seg.(*models.ConvertSegment)

	if !ok {
		return c.JSON(newAPIError(errors.Wrapf(models.ErrUnknown, "Received unknown segment type `%s`", seg.GetType())))
	}

	return c.JSON(http.StatusOK, buildConvertSegment(convSeg))
}

// GetSegmentByID _
// // GET /segments/{id} | Segment
func (dh *DealerHandler) GetSegmentByID(c echo.Context, id dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	seg, err := dh.dealer.GetSegmentByID(c.Request().Context(), author, string(id))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	convSeg, ok := seg.(*models.ConvertSegment)

	if !ok {
		return c.JSON(newAPIError(models.ErrUnknownType))
	}

	return c.JSON(http.StatusOK, buildConvertSegment(convSeg))
}

// GetSegmentsByOrderID _
// (GET /orders/{orderID}/segments)
func (dh *DealerHandler) GetSegmentsByOrderID(c echo.Context, orderID dSchema.OrderIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	segs, err := dh.dealer.GetSegmentsByOrderID(
		c.Request().Context(),
		author,
		string(orderID),
		models.EmptySegmentFilters())

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	responseSegments := make([]*dSchema.ConvertSegment, 0)

	for _, seg := range segs {
		convSeg, ok := seg.(*models.ConvertSegment)

		if !ok {
			log.Printf("Found unknown segment type (type: `%s`, id: `%s`)\n", seg.GetType(), seg.GetID())
			continue
		}

		responseSegments = append(responseSegments, buildConvertSegment(convSeg))
	}

	return c.JSON(http.StatusOK, responseSegments)
}

// AcceptSegment _
// (POST /segments/{segmentID}/actions/accept)
func (dh *DealerHandler) AcceptSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.AcceptSegment(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// CancelSegment _
// (POST /segments/{segmentID}/actions/cancel)
func (dh *DealerHandler) CancelSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	params := &dSchema.CancellationReason{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	err := dh.dealer.CancelSegment(c.Request().Context(), author, string(segmentID), params.Reason)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// AllocateSegment _
// // POST /segments
func (dh *DealerHandler) AllocateSegment(c echo.Context) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	params := &models.ConvertDealerRequest{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	seg, err := dh.dealer.AllocateSegment(c.Request().Context(), author, params)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	convSeg, ok := seg.(*models.ConvertSegment)

	if !ok {
		return c.JSON(newAPIError(errors.Wrapf(models.ErrUnknownType, "Received `%s`", seg.GetType())))
	}

	response := buildConvertSegment(convSeg)

	return c.JSON(http.StatusOK, response)
}

// FailSegment _
// (POST /segments/{id}/actions/fail)
func (dh *DealerHandler) FailSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	params := dSchema.FailureInput{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	err := dh.dealer.FailSegment(c.Request().Context(), author, string(segmentID), errors.New(params.Failure))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// FinishSegment _
// (POST /segments/{id}/actions/finish)
func (dh *DealerHandler) FinishSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.FinishSegment(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// PublishSegment _
func (dh *DealerHandler) PublishSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.PublishSegment(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// RepublishSegment _
func (dh *DealerHandler) RepublishSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.RepublishSegment(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// QuitSegment _
// (POST /segments/{id}/actions/quit)
func (dh *DealerHandler) QuitSegment(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.QuitSegment(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// SearchSegments _
// (GET /segments)
func (dh *DealerHandler) SearchSegments(ctx echo.Context) error {
	panic("not implemented") // TODO:
}

// GetInputStorageClaim _
// (GET /segments/{id}/input_storage_claim)
func (dh *DealerHandler) GetInputStorageClaim(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.GetInputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.StorageClaim{Url: storageClaim.GetURL()})
}

// AllocateInputStorageClaim _
// (POST /segments/{id}/input_storage_claim)
func (dh *DealerHandler) AllocateInputStorageClaim(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.AllocateInputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.StorageClaim{Url: storageClaim.GetURL()})
}

// GetQueuedSegmentsCount _
func (dh *DealerHandler) GetQueuedSegmentsCount(c echo.Context) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	count, err := dh.dealer.GetQueuedSegmentsCount(c.Request().Context(), author)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.Count{Count: count})
}

// NotifyProcess _
// (POST /segments/{id}/notifications/process)
func (dh *DealerHandler) NotifyProcess(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	params := &dSchema.ProgressInput{}

	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	err := dh.dealer.NotifyProcess(c.Request().Context(), author, string(segmentID), models.NewPercentProgress(float64(params.Progress)))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}

// GetOutputStorageClaim _
// (GET /segments/{id}/output_storage_claim)
func (dh *DealerHandler) GetOutputStorageClaim(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.GetOutputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.StorageClaim{Url: storageClaim.GetURL()})
}

// AllocateOutputStorageClaim _
// (POST /segments/{id}/output_storage_claim)
func (dh *DealerHandler) AllocateOutputStorageClaim(c echo.Context, segmentID dSchema.SegmentIDParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.AllocateOutputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(http.StatusOK, &dSchema.StorageClaim{Url: storageClaim.GetURL()})
}
