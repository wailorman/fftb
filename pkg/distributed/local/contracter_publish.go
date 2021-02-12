package local

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/segm"
)

// PublishOrder _
func (contracter *ContracterInstance) PublishOrder(fctx context.Context, modOrder models.IOrder) error {
	var err error

	// TODO: cancel order & segments on failure

	convOrder, ok := modOrder.(*models.ConvertOrder)

	if !ok {
		return models.ErrUnknownOrderType
	}

	slices, err := contracter.SliceConvertOrder(fctx, convOrder)

	if err != nil {
		return errors.Wrap(err, "Slicing convert order")
	}

	dSegments := make([]*models.ConvertSegment, 0)

	muxer := strings.Trim(convOrder.InFile.Extension(), ".")

	for i := range slices {
		dealerReq := &models.ConvertDealerRequest{
			Type:          models.ConvertV1Type,
			Identity:      uuid.New().String(),
			OrderIdentity: convOrder.Identity,
			Params:        convOrder.Params,
			Muxer:         muxer,
			Author:        contracter.publisher,
		}

		dealerSegment, err := contracter.dealer.AllocateSegment(dealerReq)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Allocating dealer segment #%d", i))
			// order.Failed(errObj)
			return errObj
		}

		dealerConvertSegment := dealerSegment.(*models.ConvertSegment)

		dSegments = append(dSegments, dealerConvertSegment)
		convOrder.SegmentIDs = append(convOrder.SegmentIDs, dealerSegment.GetID())
	}

	for i, slice := range slices {
		dSeg := dSegments[i]
		// seg := segs[i]
		claim, err := contracter.dealer.AllocateInputStorageClaim(contracter.publisher, dSeg.Identity)

		if err != nil {
			errObj := errors.Wrap(
				err,
				fmt.Sprintf("Allocating storage claim for dealer segment #%d (%s)", i, dSeg.GetID()),
			)

			// order.Failed(errObj)
			// TODO: cancel dealer task
			return errObj
		}

		claimWriter, err := claim.GetWriter()

		if err != nil {
			errObj := errors.Wrap(err, "Getting storage claim writer")
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return errObj
		}

		segmentReader, err := slice.File.ReadContent()

		if err != nil {
			errObj := errors.Wrap(err, "Getting segment file reader")
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return errObj
		}

		// segmentProgressReader := progress.NewReader(segmentReader)

		// // TODO: NotifyRawUpload

		// _, err = io.Copy(claimWriter, segmentProgressReader)

		_, err = io.Copy(claimWriter, segmentReader)

		if err != nil {
			return errors.Wrap(err, "Uploading segment to storage")
		}

		err = slice.File.Remove()

		if err != nil {
			return errors.Wrap(err, "Destroying slice")
		}
	}

	for _, dSeg := range dSegments {
		err := contracter.dealer.PublishSegment(contracter.publisher, dSeg.Identity)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Publishing segment %s", dSeg.Identity))
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return errObj
		}
	}

	err = contracter.registry.PersistOrder(convOrder)

	if err != nil {
		return errors.Wrap(err, "Persisting order")
	}

	return nil
}

// SliceConvertOrder _
func (contracter *ContracterInstance) SliceConvertOrder(fctx context.Context, convOrder *models.ConvertOrder) ([]*segm.Segment, error) {
	segmenter := segm.NewSliceOperation()
	segmenter.Init(segm.SliceRequest{
		InFile:         convOrder.InFile,
		KeepTimestamps: false,
		OutPath:        contracter.tempPath,
		SegmentSec:     DefaultSegmentSize,
	})

	reqSegs := make([]*segm.Segment, 0)

	sFinished, sProgress, sSegments, sFailed := segmenter.Run()

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
