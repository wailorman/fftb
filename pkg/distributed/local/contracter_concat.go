package local

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/segm"
	"golang.org/x/sync/errgroup"
)

// TODO: handle cancelled/failed

// PickOrderForConcat _
func (c *ContracterInstance) PickOrderForConcat(fctx context.Context) (models.IOrder, error) {
	// TODO: lock

	order, err := c.registry.SearchOrder(fctx, func(order models.IOrder) bool {
		segments, err := c.dealer.GetSegmentsByOrderID(fctx, order.GetID(), models.EmptySegmentFilters())

		if err != nil {
			dlog.WithOrder(c.logger, order).
				WithError(err).
				Warn("Failed to get order segments")

			return false
		}

		return order.GetCanConcat(segments)
	})

	if err != nil {
		return nil, errors.Wrap(err, "Searching for finished order")
	}

	return order, nil
}

// ConcatOrder _
func (c *ContracterInstance) ConcatOrder(fctx context.Context, order models.IOrder) error {
	// TODO: concat in local function

	logger := dlog.WithOrder(c.logger, order)

	convOrder, ok := order.(*models.ConvertOrder)

	if !ok {
		return errors.Wrapf(models.ErrUnknownOrderType, "Received `%s`", order.GetType())
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		segments, err := c.dealer.GetSegmentsByOrderID(fctx, order.GetID(), models.EmptySegmentFilters())

		if err != nil {
			return errors.Wrap(err, "Getting order segments")
		}

		slices, err := downloadSegments(c.ctx, c.dealer, c.publisher, c.storageClient, segments)

		if err != nil {
			return errors.Wrap(err, "Downloading segments")
		}

		concatOperation := segm.NewConcatOperation(
			context.WithValue(
				c.ctx,
				ctxlog.LoggerContextKey,
				logger.WithField(dlog.KeyCallee, dlog.PrefixContracterConcatWorker),
			),
		)

		err = concatOperation.Init(segm.ConcatRequest{
			OutFile:  convOrder.OutFile,
			Segments: slices,
		})

		if err != nil {
			return errors.Wrap(err, "Initializing concat operation")
		}

		cg := new(errgroup.Group)

		cg.Go(func() error {
			cProgress, cFailures := concatOperation.Run()

			var cErr error
			for {
				select {
				case pM, ok := <-cProgress:
					if ok {
						logger.WithField(dlog.KeyPercent, pM.Percent()).
							Info("Concatenating order")
					}
				case failure, failed := <-cFailures:
					if !failed {
						if cErr != nil {
							return errors.Wrap(cErr, "Failed to concatenate")
						}

						return nil
					}

					if failure != nil {
						cErr = failure
					}
				}
			}
		})

		err = cg.Wait()

		if err != nil {
			return err
		}

		err = concatOperation.Purge()

		if err != nil {
			logger.WithError(err).
				Warn("Failed to prune concat operation files")
		}

		for _, slice := range slices {
			err = slice.File.Remove()

			if err != nil {
				logger.WithField("path", slice.File.FullPath()).
					WithField("position", slice.Position).
					WithError(err).
					Warn("Failed to remove segment file")
			}
		}

		for _, segment := range segments {
			err = c.dealer.AcceptSegment(c.publisher, segment.GetID())

			if err != nil {
				dlog.WithSegment(logger, segment).
					WithError(err).
					Warn("Problem with marking segment accepted via dealer")
			}
		}

		return nil
	})

	err := g.Wait()

	if err != nil {
		return err
	}

	convOrder.State = models.OrderStateFinished

	return c.registry.PersistOrder(convOrder)
}

func downloadSegments(
	ctx context.Context,
	dealer models.IContracterDealer,
	publisher models.IAuthor,
	storageClient models.IStorageClient,
	segments []models.ISegment) ([]*segm.Segment, error) {

	slices := make([]*segm.Segment, 0)

	for _, segment := range segments {
		outputStorageClaim, err := dealer.GetOutputStorageClaim(publisher, segment.GetID())

		if err != nil {
			return nil, errors.Wrapf(err, "Getting output storage claim for segment `%s`", segment.GetID())
		}

		sliceFile, err := storageClient.MakeLocalCopy(ctx, outputStorageClaim, nil)

		if err != nil {
			return nil, errors.Wrapf(err, "Downloading output for segment `%s`", segment.GetID())
		}

		slice := &segm.Segment{
			Position: segment.GetPosition(),
			File:     sliceFile,
		}

		slices = append(slices, slice)
	}

	return slices, nil
}
