package local

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/segm"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// DefaultSegmentSize _
const DefaultSegmentSize = 10

// ContracterInstance _
type ContracterInstance struct {
	TempPath  files.Pather
	Dealer    models.IContractDealer
	Publisher models.IAuthor
}

// ContracterParameters _
type ContracterParameters struct {
	TempPath files.Pather
	Dealer   models.IContractDealer
}

// NewContracter _
func NewContracter(params *ContracterParameters) (*ContracterInstance, error) {
	publisher, err := params.Dealer.AllocatePublisherAuthority("local")

	if err != nil {
		return nil, errors.Wrap(err, "Allocating publisher authority")
	}

	return &ContracterInstance{
		TempPath:  params.TempPath,
		Dealer:    params.Dealer,
		Publisher: publisher,
	}, nil
}

// PrepareOrder _
func (c *ContracterInstance) PrepareOrder(req models.IContracterRequest) (models.IOrder, error) {
	convertRequest, ok := req.(*models.ConvertContracterRequest)

	if !ok {
		return nil, errors.Wrap(models.ErrUnknownRequestType, fmt.Sprintf("Received request with type `%s`", req.GetType()))
	}

	order := &models.ConvertOrder{
		Identity:  uuid.New().String(),
		Params:    convertRequest.Params,
		Publisher: c.Publisher,
	}

	segs, err := splitRequestToSegments(convertRequest, c.TempPath)

	if err != nil {
		errObj := errors.Wrap(err, "Splitting to segs")
		// order.Failed(errObj)
		return nil, errObj
	}

	dSegments := make([]*models.ConvertSegment, 0)

	muxer := strings.Trim(convertRequest.InFile.Extension(), ".")

	for i := range segs {
		dealerReq := &models.ConvertDealerRequest{
			Type:          models.ConvertV1Type,
			Identity:      uuid.New().String(),
			OrderIdentity: order.Identity,
			Params:        convertRequest.Params,
			Muxer:         muxer,
			Author:        c.Publisher,
		}

		dealerSegment, err := c.Dealer.AllocateSegment(dealerReq)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Allocating dealer segment #%d", i))
			// order.Failed(errObj)
			return nil, errObj
		}

		dealerConvertSegment := dealerSegment.(*models.ConvertSegment)

		dSegments = append(dSegments, dealerConvertSegment)
	}

	order.Segments = dSegments

	for i, seg := range segs {
		dSeg := dSegments[i]
		// seg := segs[i]
		claim, err := c.Dealer.AllocateInputStorageClaim(c.Publisher, dSeg.Identity)

		if err != nil {
			errObj := errors.Wrap(
				err,
				fmt.Sprintf("Allocating storage claim for dealer segment #%d (%s)", i, dSeg.GetID()),
			)

			// order.Failed(errObj)
			// TODO: cancel dealer task
			return nil, errObj
		}

		claimWriter, err := claim.GetWriter()

		if err != nil {
			errObj := errors.Wrap(err, "Getting storage claim writer")
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return nil, errObj
		}

		segmentReader, err := seg.File.ReadContent()

		if err != nil {
			errObj := errors.Wrap(err, "Getting segment file reader")
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return nil, errObj
		}

		// segmentProgressReader := progress.NewReader(segmentReader)

		// // TODO: NotifyRawUpload

		// _, err = io.Copy(claimWriter, segmentProgressReader)

		_, err = io.Copy(claimWriter, segmentReader)

		if err != nil {
			panic(err)
		}
	}

	for _, dSeg := range dSegments {
		err := c.Dealer.PublishSegment(c.Publisher, dSeg.Identity)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Publishing segment %s", dSeg.Identity))
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return nil, errObj
		}
	}

	return order, nil
}

func splitRequestToSegments(
	req *models.ConvertContracterRequest,
	// mb *models.MessageBus,
	tmpPath files.Pather,
) ([]*segm.Segment, error) {
	segmenter := segm.New()
	segmenter.Init(segm.Request{
		InFile:         req.InFile,
		KeepTimestamps: false,
		OutPath:        tmpPath,
		SegmentSec:     DefaultSegmentSize,
	})

	reqSegs := make([]*segm.Segment, 0)

	sProgress, sSegments, sFinished, sFailed := segmenter.Start()

	for {
		select {
		case <-sProgress:
		// case pr := <-sProgress:
		// mb.Publish(pr)
		case reqSeg := <-sSegments:
			reqSegs = append(reqSegs, reqSeg)
		case fail := <-sFailed:
			return nil, fail
		case <-sFinished:
			return reqSegs, nil
		}
	}
}
