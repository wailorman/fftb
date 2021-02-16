package local

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/segm"
	"golang.org/x/sync/errgroup"
)

// TODO: handle cancelled/failed

// PickOrderForConcat _
func (c *ContracterInstance) PickOrderForConcat(fctx context.Context) (models.IOrder, error) {
	var searchErr error

	ffctx, cancel := context.WithCancel(fctx)

	order, err := c.registry.SearchOrder(ffctx, func(order models.IOrder) bool {
		if order.GetState() != models.OrderStateInProgress {
			return false
		}

		statesMap, err := c.dealer.GetSegmentsStatesByOrderID(fctx, order.GetID())

		if err != nil {
			searchErr = errors.Wrap(err, "Getting segments states from dealer")
			cancel()
			return false
		}

		allSegmentsFinished := true

		for _, state := range statesMap {
			if state != models.SegmentStateFinished {
				allSegmentsFinished = false
				break
			}
		}

		return allSegmentsFinished
	})

	if errors.Is(err, models.ErrCancelled) && searchErr != nil {
		return nil, searchErr
	}

	if err != nil {
		return nil, errors.Wrap(err, "Searching for finished order")
	}

	return order, nil
}

// ConcatOrder _
func (c *ContracterInstance) ConcatOrder(fctx context.Context, order models.IOrder) error {
	convOrder, ok := order.(*models.ConvertOrder)

	if !ok {
		return errors.Wrapf(models.ErrUnknownOrderType, "Received `%s`", order.GetType())
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		slices := make([]*segm.Segment, 0)

		dSegmentsIDs := order.GetSegmentIDs()

		for _, dSegmentID := range dSegmentsIDs {
			slice, err := c.downloadSegment(fctx, order, dSegmentID)

			if err != nil {
				return errors.Wrap(err, "Downloading segment")
			}

			slices = append(slices, slice)
		}

		concatOperation := segm.NewConcatOperation(c.ctx)

		err := concatOperation.Init(segm.ConcatRequest{
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
						c.logger.WithField(dlog.KeyOrderID, order.GetID()).
							WithField(dlog.KeyPercent, pM.Percent()).
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

		err = concatOperation.Prune()

		if err != nil {
			c.logger.WithField(dlog.KeyOrderID, order.GetID()).
				WithError(err).
				Warn("Failed to prune concat operation files")
		}

		for _, slice := range slices {
			err = slice.File.Remove()

			if err != nil {
				c.logger.WithField(dlog.KeyOrderID, order.GetID()).
					WithField("path", slice.File.FullPath()).
					WithField("position", slice.Position).
					WithError(err).
					Warn("Failed to remove segment file")
			}
		}

		return nil
	})

	return g.Wait()
}

func (c *ContracterInstance) downloadSegment(fctx context.Context, order models.IOrder, segmentID string) (*segm.Segment, error) {
	dSegment, err := c.dealer.GetSegmentByID(segmentID)

	if err != nil {
		return nil, errors.Wrap(err, "Getting segment info from dealer")
	}

	segmentsOutputDir := c.tempPath.BuildSubpath("segments_output").BuildSubpath(order.GetID())

	err = segmentsOutputDir.Create()

	if err != nil {
		return nil, errors.Wrapf(err, "Creating segments output directory `%s`", segmentsOutputDir.FullPath())
	}

	sliceFile := c.tempPath.BuildSubpath(order.GetID()).BuildFile(dSegment.GetID())

	err = sliceFile.Create()

	if err != nil {
		return nil, errors.Wrapf(err, "Creating file for segment output `%s`", sliceFile.FullPath())
	}

	outputStorageClaim, err := c.dealer.GetOutputStorageClaim(c.publisher, segmentID)

	if err != nil {
		return nil, errors.Wrap(err, "Getting output storage claim")
	}

	dg := new(errgroup.Group)

	dg.Go(func() error {
		pChan, errChan := DownloadFileFromStorageClaim(fctx, sliceFile, outputStorageClaim)

		for {
			select {
			case p := <-pChan:
				if p != nil {
					c.logger.WithField(dlog.KeyPercent, p.Percent()).
						WithField(dlog.KeyOrderID, order.GetID()).
						WithField(dlog.KeySegmentID, dSegment.GetID()).
						Info("Downloading segment output")
				}

			case failure := <-errChan:
				if failure == nil {
					return nil
				}

				return failure
			}
		}
	})

	err = dg.Wait()

	if err != nil {
		return nil, err
	}

	slice := &segm.Segment{
		Position: dSegment.GetPosition(),
		File:     sliceFile,
	}

	return slice, nil
}
