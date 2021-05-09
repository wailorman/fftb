package remote

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dealerSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// DealerAPIWrapper _
type DealerAPIWrapper interface {
	// GetAllOrdersWithResponse(ctx context.Context, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetAllOrdersResponse, error)
	// GetOrderByIDWithResponse(ctx context.Context, id dealerSchema.OrderIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetOrderByIDResponse, error)
	// GetSegmentsByOrderIDWithResponse(ctx context.Context, id dealerSchema.OrderIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetSegmentsByOrderIDResponse, error)

	AllocateAuthorityWithResponse(ctx context.Context, body dealerSchema.AllocateAuthorityJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.AllocateAuthorityResponse, error)
	CreateSessionWithResponse(ctx context.Context, body dealerSchema.CreateSessionJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.CreateSessionResponse, error)
	AllocateSegmentWithResponse(ctx context.Context, body dealerSchema.AllocateSegmentJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.AllocateSegmentResponse, error)
	FindFreeSegmentWithResponse(ctx context.Context, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.FindFreeSegmentResponse, error)
	GetSegmentByIDWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetSegmentByIDResponse, error)
	FailSegmentWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, body dealerSchema.FailSegmentJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.FailSegmentResponse, error)
	FinishSegmentWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.FinishSegmentResponse, error)
	QuitSegmentWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.QuitSegmentResponse, error)
	GetInputStorageClaimWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetInputStorageClaimResponse, error)
	AllocateInputStorageClaimWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.AllocateInputStorageClaimResponse, error)
	GetOutputStorageClaimWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.GetOutputStorageClaimResponse, error)
	AllocateOutputStorageClaimWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.AllocateOutputStorageClaimResponse, error)
	NotifyProcessWithResponse(ctx context.Context, id dealerSchema.SegmentIDParam, body dealerSchema.NotifyProcessJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.NotifyProcessResponse, error)
}

// SessionCreator _
type SessionCreator interface {
	CreateSessionWithResponse(ctx context.Context, body dealerSchema.CreateSessionJSONRequestBody, reqEditors ...dealerSchema.RequestEditorFn) (*dealerSchema.CreateSessionResponse, error)
}

// Dealer _
type Dealer struct {
	apiWrapper      DealerAPIWrapper
	sc              models.IStorageClient
	authoritySecret []byte
	// TODO: basic auth
}

// NewDealer _
func NewDealer(apiWrapper DealerAPIWrapper, sc models.IStorageClient, authoritySecret []byte) *Dealer {
	return &Dealer{
		apiWrapper:      apiWrapper,
		sc:              sc,
		authoritySecret: authoritySecret,
	}
}

func buildAllocateSegmentRequest(req models.IDealerRequest) (dealerSchema.AllocateSegmentJSONRequestBody, error) {
	if req == nil {
		return dealerSchema.AllocateSegmentJSONRequestBody{},
			models.ErrMissingRequest
	}

	convReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return dealerSchema.AllocateSegmentJSONRequestBody{},
			errors.Wrapf(models.ErrUnknownType, "Unknown request type: `%s`", req.GetType())
	}

	body := dealerSchema.AllocateSegmentJSONRequestBody{
		Type:     models.ConvertV1Type,
		Id:       convReq.Identity,
		OrderId:  convReq.OrderIdentity,
		Muxer:    convReq.Muxer,
		Position: convReq.Position,

		Params: dealerSchema.ConvertParams{
			HwAccel:          convReq.Params.HWAccel,
			KeyframeInterval: convReq.Params.KeyframeInterval,
			Preset:           convReq.Params.Preset,
			Scale:            convReq.Params.Scale,
			VideoBitRate:     convReq.Params.VideoBitRate,
			VideoCodec:       convReq.Params.VideoCodec,
			VideoQuality:     convReq.Params.VideoQuality,
		},
	}

	return body, nil
}

func toModelSegment(seg *dealerSchema.ConvertSegment) (models.ISegment, error) {
	if seg == nil {
		return nil, errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	return &models.ConvertSegment{
		Identity:      seg.Id,
		OrderIdentity: seg.OrderId,
		Type:          seg.Type,
		Muxer:         seg.Muxer,
		Position:      seg.Position,
		Params: convert.Params{
			HWAccel:          seg.Params.HwAccel,
			KeyframeInterval: seg.Params.KeyframeInterval,
			Preset:           seg.Params.Preset,
			Scale:            seg.Params.Scale,
			VideoBitRate:     seg.Params.VideoBitRate,
			VideoCodec:       seg.Params.VideoCodec,
			VideoQuality:     seg.Params.VideoQuality,
		},
	}, nil
}

func withAuthor(author models.IAuthor) dealerSchema.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", author.GetSessionKey()))
		return nil
	}
}

func withUnauthorizedRetry(ctx context.Context, ca SessionCreator, author models.IAuthor, retryable func() error) error {
	firstErr := retryable()

	if firstErr == nil {
		return nil
	}

	if errors.Is(firstErr, models.ErrInvalidSessionKey) || errors.Is(firstErr, models.ErrMissingAccessToken) {
		newSessionKey, err := createSession(ctx, ca, author.GetAuthorityKey())

		if err != nil {
			return errors.Wrap(err, "Trying to create new session")
		}

		author.SetSessionKey(newSessionKey)

		return retryable()
	}

	return firstErr
}

func createSession(ctx context.Context, ca SessionCreator, authorityKey string) (string, error) {
	body := dealerSchema.CreateSessionJSONRequestBody{
		AuthorityKey: authorityKey,
	}

	response, reqErr := ca.CreateSessionWithResponse(ctx, body)

	err := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON422)

	if err != nil {
		return "", err
	}

	if response.JSON200 == nil {
		return "", errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	return response.JSON200.Key, nil
}

// AllocatePerformerAuthority _
func (rd *Dealer) AllocatePerformerAuthority(ctx context.Context, name string) (models.IAuthor, error) {
	authorityToken, err := handlers.CreateAuthorityToken(rd.authoritySecret, name)

	if err != nil {
		return nil, err
	}

	return &models.Author{Name: name, AuthorityKey: authorityToken}, nil
}

// AllocatePublisherAuthority _
func (rd *Dealer) AllocatePublisherAuthority(ctx context.Context, name string) (models.IAuthor, error) {
	authorityToken, err := handlers.CreateAuthorityToken(rd.authoritySecret, name)

	if err != nil {
		return nil, err
	}

	return &models.Author{Name: name, AuthorityKey: authorityToken}, nil
}

// AllocateSegment _
func (rd *Dealer) AllocateSegment(
	ctx context.Context,
	publisher models.IAuthor,
	req models.IDealerRequest) (models.ISegment, error) {

	body, err := buildAllocateSegmentRequest(req)

	if err != nil {
		return nil, errors.Wrap(err, "Building allocate segment request")
	}

	var response *dealerSchema.AllocateSegmentResponse
	var reqErr error

	err = withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.AllocateSegmentWithResponse(ctx, body, withAuthor(publisher))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON422, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// GetOutputStorageClaim _
func (rd *Dealer) GetOutputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	panic("not implemented") // TODO:
}

// AllocateInputStorageClaim _
func (rd *Dealer) AllocateInputStorageClaim(ctx context.Context, publisher models.IAuthor, id string) (models.IStorageClaim, error) {
	panic("not implemented") // TODO:
}

// GetQueuedSegmentsCount _
func (rd *Dealer) GetQueuedSegmentsCount(ctx context.Context, publisher models.IAuthor) (int, error) {
	panic("not implemented") // TODO:
}

// GetSegmentsByOrderID _
func (rd *Dealer) GetSegmentsByOrderID(ctx context.Context, publisher models.IAuthor, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	panic("not implemented") // TODO:
}

// GetSegmentByID _
func (rd *Dealer) GetSegmentByID(
	ctx context.Context,
	publisher models.IAuthor,
	segmentID string) (models.ISegment, error) {

	var response *dealerSchema.GetSegmentByIDResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetSegmentByIDWithResponse(ctx, dealerSchema.SegmentIDParam(segmentID), withAuthor(publisher))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// NotifyRawUpload _
func (rd *Dealer) NotifyRawUpload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
	panic("not implemented") // TODO:
}

// NotifyResultDownload _
func (rd *Dealer) NotifyResultDownload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
	panic("not implemented") // TODO:
}

// PublishSegment _
func (rd *Dealer) PublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	panic("not implemented") // TODO:
}

// RepublishSegment _
func (rd *Dealer) RepublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	panic("not implemented") // TODO:
}

// CancelSegment _
func (rd *Dealer) CancelSegment(ctx context.Context, publisher models.IAuthor, id string, reason string) error {
	panic("not implemented") // TODO:
}

// AcceptSegment _
func (rd *Dealer) AcceptSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	panic("not implemented") // TODO:
}

// ObserveSegments _
func (rd *Dealer) ObserveSegments(ctx context.Context, wg chwg.WaitGrouper) {
	panic("not implemented") // TODO:
}

// FindFreeSegment _
func (rd *Dealer) FindFreeSegment(
	ctx context.Context,
	performer models.IAuthor) (models.ISegment, error) {

	var response *dealerSchema.FindFreeSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FindFreeSegmentWithResponse(ctx, withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// NotifyRawDownload _
func (rd *Dealer) NotifyRawDownload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	panic("not implemented") // TODO:
}

// NotifyResultUpload _
func (rd *Dealer) NotifyResultUpload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	panic("not implemented") // TODO:
}

// NotifyProcess _
func (rd *Dealer) NotifyProcess(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	panic("not implemented") // TODO:
}

// FinishSegment _
func (rd *Dealer) FinishSegment(ctx context.Context, performer models.IAuthor, id string) error {
	var response *dealerSchema.FinishSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FinishSegmentWithResponse(ctx, dealerSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// QuitSegment _
func (rd *Dealer) QuitSegment(ctx context.Context, performer models.IAuthor, id string) error {
	var response *dealerSchema.QuitSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.QuitSegmentWithResponse(ctx, dealerSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// FailSegment _
func (rd *Dealer) FailSegment(ctx context.Context, performer models.IAuthor, id string, reportedErr error) error {
	var response *dealerSchema.FailSegmentResponse
	var reqErr error

	body := dealerSchema.FailSegmentJSONRequestBody{
		Failure: reportedErr.Error(),
	}

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FailSegmentWithResponse(ctx, dealerSchema.SegmentIDParam(id), body, withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// GetInputStorageClaim _
func (rd *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, id string) (models.IStorageClaim, error) {
	var response *dealerSchema.GetInputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.GetInputStorageClaimWithResponse(ctx, dealerSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return nil, errors.Wrap(err, "Calling API")
	}

	if response.JSON200 == nil {
		return nil, errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	storageClaim, err := rd.sc.BuildStorageClaimByURL(response.JSON200.Url)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim")
	}

	return storageClaim, nil
}

// AllocateOutputStorageClaim _
func (rd *Dealer) AllocateOutputStorageClaim(ctx context.Context, performer models.IAuthor, id string) (models.IStorageClaim, error) {
	var response *dealerSchema.AllocateOutputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.AllocateOutputStorageClaimWithResponse(ctx, dealerSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return errors.Wrapf(models.ErrUnknown, "Missing response (`%s`)", reqErr)
		}

		pErr := parseError(reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return nil, errors.Wrap(err, "Calling API")
	}

	if response.JSON200 == nil {
		return nil, errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	storageClaim, err := rd.sc.BuildStorageClaimByURL(response.JSON200.Url)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim")
	}

	return storageClaim, nil
}
