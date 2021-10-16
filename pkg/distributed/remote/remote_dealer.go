package remote

import (
	"context"
	"fmt"
	"net/http"

	"log"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// DealerAPIWrapper _
type DealerAPIWrapper interface {
	// GetAllOrdersWithResponse(ctx context.Context, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetAllOrdersResponse, error)
	// GetOrderByIDWithResponse(ctx context.Context, id dSchema.OrderIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetOrderByIDResponse, error)
	// GetSegmentsByOrderIDWithResponse(ctx context.Context, id dSchema.OrderIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetSegmentsByOrderIDResponse, error)

	AllocateAuthorityWithResponse(ctx context.Context, body dSchema.AllocateAuthorityJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.AllocateAuthorityResponse, error)
	CreateSessionWithResponse(ctx context.Context, body dSchema.CreateSessionJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.CreateSessionResponse, error)
	AllocateSegmentWithResponse(ctx context.Context, body dSchema.AllocateSegmentJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.AllocateSegmentResponse, error)
	FindFreeSegmentWithResponse(ctx context.Context, reqEditors ...dSchema.RequestEditorFn) (*dSchema.FindFreeSegmentResponse, error)
	GetSegmentByIDWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetSegmentByIDResponse, error)
	FailSegmentWithResponse(ctx context.Context, id dSchema.SegmentIDParam, body dSchema.FailSegmentJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.FailSegmentResponse, error)
	FinishSegmentWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.FinishSegmentResponse, error)
	QuitSegmentWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.QuitSegmentResponse, error)
	GetInputStorageClaimWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetInputStorageClaimResponse, error)
	AllocateInputStorageClaimWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.AllocateInputStorageClaimResponse, error)
	GetOutputStorageClaimWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetOutputStorageClaimResponse, error)
	AllocateOutputStorageClaimWithResponse(ctx context.Context, id dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.AllocateOutputStorageClaimResponse, error)
	NotifyProcessWithResponse(ctx context.Context, id dSchema.SegmentIDParam, body dSchema.NotifyProcessJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.NotifyProcessResponse, error)
	SearchSegmentsWithResponse(ctx context.Context, reqEditors ...dSchema.RequestEditorFn) (*dSchema.SearchSegmentsResponse, error)
	GetSegmentsByOrderIDWithResponse(ctx context.Context, orderID dSchema.OrderIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetSegmentsByOrderIDResponse, error)
	AcceptSegmentWithResponse(ctx context.Context, segmentID dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.AcceptSegmentResponse, error)
	CancelSegmentWithResponse(ctx context.Context, segmentID dSchema.SegmentIDParam, body dSchema.CancelSegmentJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.CancelSegmentResponse, error)
	PublishSegmentWithResponse(ctx context.Context, segmentID dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.PublishSegmentResponse, error)
	RepublishSegmentWithResponse(ctx context.Context, segmentID dSchema.SegmentIDParam, reqEditors ...dSchema.RequestEditorFn) (*dSchema.RepublishSegmentResponse, error)
	GetQueuedSegmentsCountWithResponse(ctx context.Context, reqEditors ...dSchema.RequestEditorFn) (*dSchema.GetQueuedSegmentsCountResponse, error)
}

// SessionCreator _
type SessionCreator interface {
	CreateSessionWithResponse(ctx context.Context, body dSchema.CreateSessionJSONRequestBody, reqEditors ...dSchema.RequestEditorFn) (*dSchema.CreateSessionResponse, error)
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

func buildAllocateSegmentRequest(req models.IDealerRequest) (dSchema.AllocateSegmentJSONRequestBody, error) {
	if req == nil {
		return dSchema.AllocateSegmentJSONRequestBody{},
			models.ErrMissingRequest
	}

	convReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return dSchema.AllocateSegmentJSONRequestBody{},
			errors.Wrapf(models.ErrUnknownType, "Unknown request type: `%s`", req.GetType())
	}

	body := dSchema.AllocateSegmentJSONRequestBody{
		Type:     models.ConvertV1Type,
		Id:       convReq.Identity,
		OrderId:  convReq.OrderIdentity,
		Muxer:    convReq.Muxer,
		Position: convReq.Position,

		Params: dSchema.ConvertParams{
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

func toModelSegment(seg *dSchema.ConvertSegment) (models.ISegment, error) {
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

func withAuthor(author models.IAuthor) dSchema.RequestEditorFn {
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
	body := dSchema.CreateSessionJSONRequestBody{
		AuthorityKey: authorityKey,
	}

	response, reqErr := ca.CreateSessionWithResponse(ctx, body)

	err := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON422)

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
func (rd *Dealer) AllocateSegment(ctx context.Context, publisher models.IAuthor, req models.IDealerRequest) (models.ISegment, error) {

	body, err := buildAllocateSegmentRequest(req)

	if err != nil {
		return nil, errors.Wrap(err, "Building allocate segment request")
	}

	var response *dSchema.AllocateSegmentResponse
	var reqErr error

	err = withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.AllocateSegmentWithResponse(ctx, body, withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON422, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// GetOutputStorageClaim _
func (rd *Dealer) GetOutputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
	var response *dSchema.GetOutputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetOutputStorageClaimWithResponse(ctx, dSchema.SegmentIDParam(segmentID), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
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

// AllocateInputStorageClaim _
func (rd *Dealer) AllocateInputStorageClaim(ctx context.Context, publisher models.IAuthor, id string) (models.IStorageClaim, error) {
	var response *dSchema.AllocateInputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.AllocateInputStorageClaimWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON404, response.JSON401, response.JSON403)
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

// GetQueuedSegmentsCount _
func (rd *Dealer) GetQueuedSegmentsCount(ctx context.Context, publisher models.IAuthor) (int, error) {
	var response *dSchema.GetQueuedSegmentsCountResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetQueuedSegmentsCountWithResponse(ctx, withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return 0, errors.Wrap(err, "Calling API")
	}

	return response.JSON200.Count, nil
}

// GetSegmentsByOrderID _
func (rd *Dealer) GetSegmentsByOrderID(ctx context.Context, publisher models.IAuthor, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	var response *dSchema.GetSegmentsByOrderIDResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetSegmentsByOrderIDWithResponse(ctx, dSchema.OrderIDParam(orderID), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	segments := make([]models.ISegment, 0)

	for _, segment := range *response.JSON200 {
		mSeg, err := toModelSegment(&segment)

		if err != nil {
			return nil, errors.Wrap(err, "Parsing segment")
		}

		segments = append(segments, mSeg)
	}

	return segments, nil
}

// GetSegmentByID _
func (rd *Dealer) GetSegmentByID(ctx context.Context, publisher models.IAuthor, segmentID string) (models.ISegment, error) {

	var response *dSchema.GetSegmentByIDResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetSegmentByIDWithResponse(ctx, dSchema.SegmentIDParam(segmentID), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// NotifyRawUpload _
func (rd *Dealer) NotifyRawUpload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
	// TODO: not implemented
	log.Printf("remote.NotifyRawUpload (not implemented): #%s %.2f\n", id, p.Percent())
	return nil
}

// NotifyResultDownload _
func (rd *Dealer) NotifyResultDownload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
	// TODO: not implemented
	log.Printf("remote.NotifyResultDownload (not implemented): #%s %.2f\n", id, p.Percent())
	return nil
}

// PublishSegment _
func (rd *Dealer) PublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	var response *dSchema.PublishSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.PublishSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// RepublishSegment _
func (rd *Dealer) RepublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	var response *dSchema.RepublishSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.RepublishSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// CancelSegment _
func (rd *Dealer) CancelSegment(ctx context.Context, publisher models.IAuthor, id string, reason string) error {
	var response *dSchema.CancelSegmentResponse
	var reqErr error

	requestBody := dSchema.CancelSegmentJSONRequestBody{Reason: reason}

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.CancelSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), requestBody, withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// AcceptSegment _
func (rd *Dealer) AcceptSegment(ctx context.Context, publisher models.IAuthor, id string) error {
	var response *dSchema.AcceptSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.AcceptSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(publisher))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// ObserveSegments _
func (rd *Dealer) ObserveSegments(ctx context.Context, wg chwg.WaitGrouper) {
	// TODO: not implemented
	log.Printf("remote.ObserveSegments (not implemented) SHOULD BE REMOVED\n")
}

// FindFreeSegment _
func (rd *Dealer) FindFreeSegment(ctx context.Context, performer models.IAuthor) (models.ISegment, error) {
	var response *dSchema.FindFreeSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FindFreeSegmentWithResponse(ctx, withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// NotifyRawDownload _
func (rd *Dealer) NotifyRawDownload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	// TODO: not implemented
	log.Printf("remote.NotifyRawDownload (not implemented): #%s %.2f\n", id, p.Percent())
	return nil
}

// NotifyResultUpload _
func (rd *Dealer) NotifyResultUpload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	// TODO: not implemented
	log.Printf("remote.NotifyResultUpload (not implemented): #%s %.2f\n", id, p.Percent())
	return nil
}

// NotifyProcess _
func (rd *Dealer) NotifyProcess(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
	var response *dSchema.NotifyProcessResponse
	var reqErr error

	requestBody := dSchema.NotifyProcessJSONRequestBody{
		Progress: float32(p.Percent()),
	}

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.NotifyProcessWithResponse(ctx, dSchema.SegmentIDParam(id), requestBody, withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// FinishSegment _
func (rd *Dealer) FinishSegment(ctx context.Context, performer models.IAuthor, id string) error {
	var response *dSchema.FinishSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FinishSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// QuitSegment _
func (rd *Dealer) QuitSegment(ctx context.Context, performer models.IAuthor, id string) error {
	var response *dSchema.QuitSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.QuitSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// FailSegment _
func (rd *Dealer) FailSegment(ctx context.Context, performer models.IAuthor, id string, reportedErr error) error {
	var response *dSchema.FailSegmentResponse
	var reqErr error

	body := dSchema.FailSegmentJSONRequestBody{
		Failure: reportedErr.Error(),
	}

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FailSegmentWithResponse(ctx, dSchema.SegmentIDParam(id), body, withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(noContentBody, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
		return pErr
	})

	if err != nil {
		return err
	}

	return nil
}

// GetInputStorageClaim _
func (rd *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, id string) (models.IStorageClaim, error) {
	var response *dSchema.GetInputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.GetInputStorageClaimWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
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
	var response *dSchema.AllocateOutputStorageClaimResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.AllocateOutputStorageClaimWithResponse(ctx, dSchema.SegmentIDParam(id), withAuthor(performer))

		if response == nil {
			return buildEmptyResponseError(reqErr, response.HTTPResponse, response.Body)
		}

		pErr := parseError(response.JSON200, reqErr, response.HTTPResponse, response.Body, response.JSON401, response.JSON403, response.JSON404)
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
