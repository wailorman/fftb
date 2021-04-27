package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/schema"
)

// DealerHandler _
type DealerHandler struct {
	ctx             context.Context
	dealer          models.IDealer
	authoritySecret []byte
	sessionSecret   []byte
}

func buildConvertSegment(convSeg *models.ConvertSegment) *schema.ConvertSegment {
	return &schema.ConvertSegment{
		Id:       convSeg.Identity,
		OrderId:  convSeg.OrderIdentity,
		Type:     convSeg.Type,
		Muxer:    convSeg.Muxer,
		Position: convSeg.Position,
		Params: schema.ConvertParams{
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

// NewDealerHandler _
func NewDealerHandler(
	ctx context.Context,
	dealer models.IDealer,
	authoritySecret []byte,
	sessionSecret []byte) *DealerHandler {

	return &DealerHandler{
		ctx:             ctx,
		dealer:          dealer,
		authoritySecret: authoritySecret,
		sessionSecret:   sessionSecret,
	}
}

// AllocateAuthority _
// POST /authorities
func (dh *DealerHandler) AllocateAuthority(c echo.Context) error {
	params := &schema.AuthorityInput{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	key, err := CreateAuthorityToken(dh.authoritySecret, params.Name)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(200, &schema.Authority{Key: key})
}

// CreateSession _
// POST /sessions
func (dh *DealerHandler) CreateSession(c echo.Context) error {
	params := &schema.SessionInput{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(newAPIError(err))
	}

	key, err := CreateSessionToken(dh.authoritySecret, dh.sessionSecret, params.AuthorityKey)

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(200, &schema.Session{Key: key})
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

	return c.JSON(200, buildConvertSegment(convSeg))
}

// GetSegmentByID _
// // GET /segments/{id} | Segment
func (dh *DealerHandler) GetSegmentByID(c echo.Context, id schema.SegmentIdParam) error {
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

	return c.JSON(200, buildConvertSegment(convSeg))
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

	return c.JSON(200, response)
}

// FailSegment _
// (POST /segments/{id}/actions/fail)
func (dh *DealerHandler) FailSegment(c echo.Context, segmentID schema.SegmentIdParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	params := schema.FailureInput{}
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
func (dh *DealerHandler) FinishSegment(c echo.Context, segmentID schema.SegmentIdParam) error {
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

// QuitSegment _
// (POST /segments/{id}/actions/quit)
func (dh *DealerHandler) QuitSegment(c echo.Context, segmentID schema.SegmentIdParam) error {
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

// GetInputStorageClaim _
// (GET /segments/{id}/input_storage_claim)
func (dh *DealerHandler) GetInputStorageClaim(c echo.Context, segmentID schema.SegmentIdParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.GetInputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(200, &schema.StorageClaim{Url: storageClaim.GetURL()})
}

// AllocateInputStorageClaim _
// (POST /segments/{id}/input_storage_claim)
func (dh *DealerHandler) AllocateInputStorageClaim(c echo.Context, segmentID schema.SegmentIdParam) error {
	panic("not implemented")
}

// NotifyProcess _
// (POST /segments/{id}/notifications/process)
func (dh *DealerHandler) NotifyProcess(c echo.Context, segmentID schema.SegmentIdParam) error {
	panic("not implemented")
}

// GetOutputStorageClaim _
// (GET /segments/{id}/output_storage_claim)
func (dh *DealerHandler) GetOutputStorageClaim(c echo.Context, segmentID schema.SegmentIdParam) error {
	panic("not implemented")
}

// AllocateOutputStorageClaim _
// (POST /segments/{id}/output_storage_claim)
func (dh *DealerHandler) AllocateOutputStorageClaim(c echo.Context, segmentID schema.SegmentIdParam) error {
	author := extractAuthor(c)

	if author == nil {
		return c.JSON(newAPIError(models.ErrMissingAuthor))
	}

	storageClaim, err := dh.dealer.AllocateOutputStorageClaim(c.Request().Context(), author, string(segmentID))

	if err != nil {
		return c.JSON(newAPIError(err))
	}

	return c.JSON(200, &schema.StorageClaim{Url: storageClaim.GetURL()})
}
