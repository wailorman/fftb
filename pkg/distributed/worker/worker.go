package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/throttle"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/minfo"
)

// TODO: do not pass regular context to cancellation code

// FreeTaskDelay _
const FreeTaskDelay = time.Duration(3) * time.Second

// FreeTaskTimeout _
const FreeTaskTimeout = time.Duration(3) * time.Second

// Instance _
type Instance struct {
	ctx           context.Context
	tmpPath       files.Pather
	dealer        models.IWorkerDealer
	logger        logrus.FieldLogger
	performer     models.IAuthor
	storageClient models.IStorageClient
	wg            *chwg.ChannelledWaitGroup
}

// segmentIO _
type segmentIO struct {
	inputClaim  models.IStorageClaim
	outputClaim models.IStorageClaim
	inputFile   files.Filer
	outputFile  files.Filer
}

// NewWorker _
func NewWorker(ctx context.Context, tmpPath files.Pather, dealer models.IWorkerDealer, storageClient models.IStorageClient) (*Instance, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.worker"); logger == nil {
		logger = ctxlog.New("fftb.worker")
	}

	performer, err := dealer.AllocatePerformerAuthority(ctx, uuid.New().String())

	if err != nil {
		return nil, errors.Wrap(err, "Obtaining performer authority")
	}

	return &Instance{
		ctx:           ctx,
		tmpPath:       tmpPath,
		dealer:        dealer,
		logger:        logger,
		performer:     performer,
		storageClient: storageClient,
		wg:            chwg.New(),
	}, nil
}

// Start _
func (w *Instance) Start() {
	w.wg.Add(1)

	go func() {
		for {
			freeSegment, err := w.dealer.FindFreeSegment(w.ctx, w.performer)

			if err != nil {
				if errors.Is(err, models.ErrNotFound) {
					w.logger.Debug("Free segment not found")
				} else {
					w.logger.WithError(err).Warn("Searching free segment error")
				}

				time.Sleep(FreeTaskDelay)
				continue
			}

			logger := dlog.WithSegment(w.logger, freeSegment).
				WithField(dlog.KeyPerformer, w.performer.GetName())

			logger.Info("Found free segment")

			err = proceedSegment(
				w.ctx,
				w.wg,
				logger,
				w.performer,
				w.dealer,
				w.storageClient,
				w.tmpPath,
				freeSegment)

			if err != nil {
				logger.WithError(err).Warn("Processing segment error")
			}
		}
	}()
}

func proceedSegment(
	ctx context.Context,
	wg chwg.WaitGrouper,
	logger *logrus.Entry,
	performer models.IAuthor,
	dealer models.IWorkerDealer,
	storageClient models.IStorageClient,
	tmpPath files.Pather,
	freeSegment models.ISegment) error {

	wg.Add(1)
	defer wg.Done()

	fail := func(err error) error {
		return failSegment(ctx, logger, wg, performer, dealer, freeSegment, err)
	}

	convertSegment, ok := freeSegment.(*models.ConvertSegment)

	if !ok {
		return fail(errors.Wrapf(models.ErrUnknownType, "Received type `%s`", freeSegment.GetType()))
	}

	sio, err := prepareSegmentIO(ctx, performer, dealer, storageClient, convertSegment, tmpPath)

	if err != nil {
		return fail(errors.Wrap(err, "Preparing segment IO"))
	}

	task := convert.Task{
		InFile:  sio.inputFile.FullPath(),
		OutFile: sio.outputFile.FullPath(),
		Params:  convertSegment.Params,
	}

	converter := convert.NewConverter(
		context.WithValue(
			ctx,
			ctxlog.LoggerContextKey,
			logger.WithField(dlog.KeyCallee, dlog.PrefixWorker),
		),
		minfo.New(),
	)

	throttled := throttle.New(2000 * time.Millisecond)

	progressChan, errChan := converter.Convert(task)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Terminating worker thread")

			// <-converter.Closed()

			if err = dealer.QuitSegment(context.Background(), performer, freeSegment.GetID()); err != nil {
				logger.WithError(err).Warn("Problem with quiting segment")
			}

			if err = storageClient.RemoveLocalCopy(ctx, sio.inputClaim); err != nil {
				logger.WithError(err).Warn("Problem with removing input local copy file")
			}

			if err = sio.outputFile.Remove(); err != nil {
				logger.WithError(err).Warn("Problem with removing output file")
			}

			return ctx.Err()

		case pmsg, ok := <-progressChan:
			if ok {
				throttled(func() {
					modProgress := makeProgresserFromConvert(pmsg)

					if err = dealer.NotifyProcess(ctx, performer, freeSegment.GetID(), modProgress); err != nil {
						logger.WithError(err).Warn("Problem with notifying process")
					}

					dlog.SegmentProgress(logger, freeSegment, modProgress)
				})
			}

		case cErr, failed := <-errChan:
			if !failed {
				err := moveOutput(ctx, storageClient, sio.outputFile, sio.outputClaim)

				if err != nil {
					return fail(errors.Wrap(err, "Uploading output & removing local copy"))
				}

				err = storageClient.RemoveLocalCopy(ctx, sio.inputClaim)

				if err != nil {
					logger.WithError(err).Warn("Problem with removing input file")
				}

				logger.Info("Segment is done")

				if err = dealer.FinishSegment(context.Background(), performer, freeSegment.GetID()); err != nil {
					return fail(errors.Wrap(err, "Sending segment finish report"))
				}

				return nil
			}

			return fail(errors.Wrap(cErr, "Error processing convert task"))
		}
	}
}

func moveOutput(
	ctx context.Context,
	storageClient models.IStorageClient,
	outputFile files.Filer,
	outputClaim models.IStorageClaim) error {

	err := storageClient.MoveFileToStorageClaim(ctx, outputFile, outputClaim, nil)

	if err != nil {
		return errors.Wrap(err, "Moving (uploading) output file to storage claim")
	}

	return nil
}

func failSegment(
	ctx context.Context,
	logger logrus.FieldLogger,
	wg chwg.WaitGrouper,
	performer models.IAuthor,
	dealer models.IWorkerDealer,
	segment models.ISegment,
	err error) error {

	if errors.Is(err, context.Canceled) {
		return err
	}

	wg.Add(1)

	dErr := dealer.FailSegment(context.Background(), performer, segment.GetID(), err)

	if dErr != nil {
		logger.WithError(err).
			Warn("Failed to report segment failure")
	}

	wg.Done()

	return err
}

func prepareSegmentIO(
	ctx context.Context,
	performer models.IAuthor,
	dealer models.IWorkerDealer,
	storageClient models.IStorageClient,
	convSegment *models.ConvertSegment,
	tmpPath files.Pather) (*segmentIO, error) {

	inputClaim, err := dealer.GetInputStorageClaim(ctx, performer, convSegment.GetID())

	if err != nil {
		return nil, errors.Wrap(err, "Getting input storage claim")
	}

	outputClaim, err := dealer.AllocateOutputStorageClaim(ctx, performer, convSegment.GetID())

	if err != nil {
		return nil, errors.Wrap(err, "Getting output storage claim")
	}

	inputFile, err := storageClient.MakeLocalCopy(ctx, inputClaim, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Downloading local copy of input storage claim")
	}

	outputFile := tmpPath.BuildFile(fmt.Sprintf("%s.%s", outputClaim.GetID(), convSegment.Muxer))

	return &segmentIO{
		inputClaim:  inputClaim,
		outputClaim: outputClaim,
		inputFile:   inputFile,
		outputFile:  outputFile,
	}, nil
}

// Closed _
func (w *Instance) Closed() <-chan struct{} {
	return w.wg.Closed()
}
