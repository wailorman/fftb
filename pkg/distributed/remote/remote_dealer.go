package remote

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// APIWrapper _
type APIWrapper interface {
	AllocateSegmentWithResponse(ctx context.Context, body AllocateSegmentJSONRequestBody, reqEditors ...RequestEditorFn) (*AllocateSegmentResponse, error)
	GetSegmentByIDWithResponse(ctx context.Context, id SegmentIdParam, reqEditors ...RequestEditorFn) (*GetSegmentByIDResponse, error)
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

func buildAllocateSegmentRequest(req models.IDealerRequest) (AllocateSegmentJSONRequestBody, error) {
	if req == nil {
		return AllocateSegmentJSONRequestBody{},
			models.ErrMissingRequest
	}

	convReq, ok := req.(*models.ConvertDealerRequest)

	if !ok {
		return AllocateSegmentJSONRequestBody{},
			errors.Wrapf(models.ErrUnknownType, "Unknown request type: `%s`", req.GetType())
	}

	body := AllocateSegmentJSONRequestBody{
		Type:     models.ConvertV1Type,
		Id:       convReq.Identity,
		OrderId:  convReq.OrderIdentity,
		Muxer:    &convReq.Muxer,
		Position: &convReq.Position,

		Params: &ConvertParams{
			HwAccel:          &convReq.Params.HWAccel,
			KeyframeInterval: &convReq.Params.KeyframeInterval,
			Preset:           &convReq.Params.Preset,
			Scale:            &convReq.Params.Scale,
			VideoBitRate:     &convReq.Params.VideoBitRate,
			VideoCodec:       &convReq.Params.VideoCodec,
			VideoQuality:     &convReq.Params.VideoQuality,
		},
	}

	return body, nil
}

func toModelSegment(seg *Segment) (models.ISegment, error) {
	if seg == nil {
		return nil, errors.Wrap(models.ErrUnknown, "Missing success response")
	}

	return &models.ConvertSegment{
		Identity: seg.Id,
	}, nil
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

	response, reqErr := rd.apiWrapper.AllocateSegmentWithResponse(ctx, body)

	err = parseError(reqErr, response.HTTPResponse, response.JSON422)

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
	response, reqErr := rd.apiWrapper.GetSegmentByIDWithResponse(ctx, SegmentIdParam(segmentID))

	err := parseError(reqErr, response.HTTPResponse, response.JSON404)

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

// // FindFreeSegment _
// func (rd *Dealer) FindFreeSegment(ctx context.Context, performer models.IAuthor) (models.ISegment, error) {
// 	panic("not implemented")
// }

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
