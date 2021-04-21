package remote

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/schema"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// APIWrapper _
type APIWrapper interface {
	AllocateSegmentWithResponse(ctx context.Context, body schema.AllocateSegmentJSONRequestBody, reqEditors ...schema.RequestEditorFn) (*schema.AllocateSegmentResponse, error)
	GetSegmentByIDWithResponse(ctx context.Context, id schema.SegmentIdParam, reqEditors ...schema.RequestEditorFn) (*schema.GetSegmentByIDResponse, error)
	AllocateAuthorityWithResponse(ctx context.Context, body schema.AllocateAuthorityJSONRequestBody, reqEditors ...schema.RequestEditorFn) (*schema.AllocateAuthorityResponse, error)
	CreateSessionWithResponse(ctx context.Context, body schema.CreateSessionJSONRequestBody, reqEditors ...schema.RequestEditorFn) (*schema.CreateSessionResponse, error)
	FindFreeSegmentWithResponse(ctx context.Context, reqEditors ...schema.RequestEditorFn) (*schema.FindFreeSegmentResponse, error)
}

// SessionCreator _
type SessionCreator interface {
	CreateSessionWithResponse(ctx context.Context, body schema.CreateSessionJSONRequestBody, reqEditors ...schema.RequestEditorFn) (*schema.CreateSessionResponse, error)
}

// Dealer _
type Dealer struct {
	apiWrapper APIWrapper
}

// NewDealer _
func NewDealer(apiWrapper APIWrapper) *Dealer {
	return &Dealer{
		apiWrapper: apiWrapper,
	}
}

// // AllocatePublisherAuthority _
// func (rd *Dealer) AllocatePublisherAuthority(ctx context.Context, name string) (models.IAuthor, error) {
// 	panic("not implemented")
// }

func buildAllocateSegmentRequest(req models.IDealerRequest) (schema.AllocateSegmentJSONRequestBody, error) {
	if req == nil {
		return schema.AllocateSegmentJSONRequestBody{},
			models.ErrMissingRequest
	}

	convReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return schema.AllocateSegmentJSONRequestBody{},
			errors.Wrapf(models.ErrUnknownType, "Unknown request type: `%s`", req.GetType())
	}

	body := schema.AllocateSegmentJSONRequestBody{
		Type:     models.ConvertV1Type,
		Id:       convReq.Identity,
		OrderId:  convReq.OrderIdentity,
		Muxer:    convReq.Muxer,
		Position: convReq.Position,

		Params: schema.ConvertParams{
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

func toModelSegment(seg *schema.ConvertSegment) (models.ISegment, error) {
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

func withAuthor(author models.IAuthor) schema.RequestEditorFn {
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
	body := schema.CreateSessionJSONRequestBody{
		AuthorityKey: authorityKey,
	}

	response, reqErr := ca.CreateSessionWithResponse(ctx, body)

	err := parseError(reqErr, response.HTTPResponse, response.JSON422)

	if err != nil {
		return "", err
	}

	if response.JSON200 == nil {
		return "", errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	return response.JSON200.Key, nil
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

	var response *schema.AllocateSegmentResponse
	var reqErr error

	err = withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.AllocateSegmentWithResponse(ctx, body, withAuthor(publisher))
		pErr := parseError(reqErr, response.HTTPResponse, response.JSON422, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// // GetOutputStorageClaim _
// func (rd *Dealer) GetOutputStorageClaim(ctx context.Context, publisher models.IAuthor, segmentID string) (models.IStorageClaim, error) {
// 	panic("not implemented")
// }

// // AllocateInputStorageClaim _
// func (rd *Dealer) AllocateInputStorageClaim(ctx context.Context, publisher models.IAuthor, id string) (models.IStorageClaim, error) {
// 	panic("not implemented")
// }

// // GetQueuedSegmentsCount _
// func (rd *Dealer) GetQueuedSegmentsCount(ctx context.Context, publisher models.IAuthor) (int, error) {
// 	panic("not implemented")
// }

// // GetSegmentsByOrderID _
// func (rd *Dealer) GetSegmentsByOrderID(ctx context.Context, publisher models.IAuthor, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
// 	panic("not implemented")
// }

// GetSegmentByID _
func (rd *Dealer) GetSegmentByID(ctx context.Context, publisher models.IAuthor, segmentID string) (models.ISegment, error) {
	var response *schema.GetSegmentByIDResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, publisher, func() error {
		response, reqErr = rd.apiWrapper.GetSegmentByIDWithResponse(ctx, schema.SegmentIdParam(segmentID), withAuthor(publisher))
		pErr := parseError(reqErr, response.HTTPResponse, response.JSON404, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// // NotifyRawUpload _
// func (rd *Dealer) NotifyRawUpload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
// 	panic("not implemented")
// }

// // NotifyResultDownload _
// func (rd *Dealer) NotifyResultDownload(ctx context.Context, publisher models.IAuthor, id string, p models.Progresser) error {
// 	panic("not implemented")
// }

// // PublishSegment _
// func (rd *Dealer) PublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
// 	panic("not implemented")
// }

// // RepublishSegment _
// func (rd *Dealer) RepublishSegment(ctx context.Context, publisher models.IAuthor, id string) error {
// 	panic("not implemented")
// }

// // CancelSegment _
// func (rd *Dealer) CancelSegment(ctx context.Context, publisher models.IAuthor, id string, reason string) error {
// 	panic("not implemented")
// }

// // AcceptSegment _
// func (rd *Dealer) AcceptSegment(ctx context.Context, publisher models.IAuthor, id string) error {
// 	panic("not implemented")
// }

// // ObserveSegments _
// func (rd *Dealer) ObserveSegments(ctx context.Context, wg chwg.WaitGrouper) {
// 	panic("not implemented")
// }

// // AllocatePerformerAuthority _
// func (rd *Dealer) AllocatePerformerAuthority(ctx context.Context, name string) (models.IAuthor, error) {
// 	panic("not implemented")
// }

// FindFreeSegment _
func (rd *Dealer) FindFreeSegment(ctx context.Context, performer models.IAuthor) (models.ISegment, error) {
	var response *schema.FindFreeSegmentResponse
	var reqErr error

	err := withUnauthorizedRetry(ctx, rd.apiWrapper, performer, func() error {
		response, reqErr = rd.apiWrapper.FindFreeSegmentWithResponse(ctx, withAuthor(performer))
		pErr := parseError(reqErr, response.HTTPResponse, response.JSON404, response.JSON401)
		return pErr
	})

	if err != nil {
		return nil, err
	}

	return toModelSegment(response.JSON200)
}

// // NotifyRawDownload _
// func (rd *Dealer) NotifyRawDownload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
// 	panic("not implemented")
// }

// // NotifyResultUpload _
// func (rd *Dealer) NotifyResultUpload(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
// 	panic("not implemented")
// }

// // NotifyProcess _
// func (rd *Dealer) NotifyProcess(ctx context.Context, performer models.IAuthor, id string, p models.Progresser) error {
// 	panic("not implemented")
// }

// // FinishSegment _
// func (rd *Dealer) FinishSegment(ctx context.Context, performer models.IAuthor, id string) error {
// 	panic("not implemented")
// }

// // QuitSegment _
// func (rd *Dealer) QuitSegment(ctx context.Context, performer models.IAuthor, id string) error {
// 	panic("not implemented")
// }

// // FailSegment _
// func (rd *Dealer) FailSegment(ctx context.Context, performer models.IAuthor, id string, err error) error {
// 	panic("not implemented")
// }

// // GetInputStorageClaim _
// func (rd *Dealer) GetInputStorageClaim(ctx context.Context, performer models.IAuthor, segmentID string) (models.IStorageClaim, error) {
// 	panic("not implemented")
// }

// // AllocateOutputStorageClaim _
// func (rd *Dealer) AllocateOutputStorageClaim(ctx context.Context, performer models.IAuthor, id string) (models.IStorageClaim, error) {
// 	panic("not implemented")
// }
