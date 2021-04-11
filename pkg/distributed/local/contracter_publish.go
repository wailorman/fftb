package local

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/segm"
)

// PublishOrder _
func (contracter *ContracterInstance) publishOrder(fctx context.Context, modOrder models.IOrder) error {
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

	for i, slice := range slices {
		dealerReq := &models.ConvertDealerRequest{
			Type:          models.ConvertV1Type,
			Identity:      uuid.New().String(),
			OrderIdentity: convOrder.Identity,
			Params:        convOrder.Params,
			Muxer:         muxer,
			Position:      slice.Position,
		}

		dealerSegment, err := contracter.dealer.AllocateSegment(fctx, contracter.publisher, dealerReq)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Allocating dealer segment #%d", i))
			// order.Failed(errObj)
			return errObj
		}

		dealerConvertSegment := dealerSegment.(*models.ConvertSegment)

		dSegments = append(dSegments, dealerConvertSegment)
	}

	for i, slice := range slices {
		dSeg := dSegments[i]
		// seg := segs[i]
		claim, err := contracter.dealer.AllocateInputStorageClaim(fctx, contracter.publisher, dSeg.Identity)

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
		claimWriter.Close()
		segmentReader.Close()

		if err != nil {
			return errors.Wrap(err, "Uploading segment to storage")
		}

		err = slice.File.Remove()

		if err != nil {
			return errors.Wrap(err, "Destroying slice")
		}
	}

	for _, dSeg := range dSegments {
		err := contracter.dealer.PublishSegment(fctx, contracter.publisher, dSeg.Identity)

		if err != nil {
			errObj := errors.Wrap(err, fmt.Sprintf("Publishing segment %s", dSeg.Identity))
			// order.Failed(errObj)
			// TODO: cancel dealer task
			return errObj
		}
	}

	convOrder.State = models.OrderStateInProgress
	err = contracter.registry.PersistOrder(fctx, convOrder)

	if err != nil {
		return errors.Wrap(err, "Persisting order")
	}

	return nil
}

// SliceConvertOrder _
func (contracter *ContracterInstance) SliceConvertOrder(fctx context.Context, convOrder *models.ConvertOrder) ([]*segm.Segment, error) {
	llog := dlog.WithOrder(contracter.logger, convOrder)

	llog.Info("Slicing order")

	segmenter := segm.NewSliceOperation(contracter.ctx)
	err := segmenter.Init(segm.SliceRequest{
		InFile:         convOrder.InFile,
		KeepTimestamps: false,
		OutPath:        contracter.tempPath,
		SegmentSec:     DefaultSegmentSize,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Initializing slice operation")
	}

	reqSegs := make([]*segm.Segment, 0)

	sProgress, sSegments, sFailed := segmenter.Run()

	for {
		select {
		case p, ok := <-sProgress:
			if ok {
				llog.
					WithField(dlog.KeyPercent, p.Percent()).
					Debug("Slicing order progress")
			}

		case reqSeg, ok := <-sSegments:
			if ok {
				reqSegs = append(reqSegs, reqSeg)
			}

		case failure, failed := <-sFailed:
			if !failed {
				return reqSegs, nil
			}

			return nil, failure
		}
	}
}
