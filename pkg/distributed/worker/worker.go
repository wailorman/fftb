package worker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/machinebox/progress"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/throttle"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/info"
	// "github.com/bep/debounce"
	// "github.com/machinebox/progress"
)

// FreeTaskDelay _
const FreeTaskDelay = time.Duration(3) * time.Second

// FreeTaskTimeout _
const FreeTaskTimeout = time.Duration(3) * time.Second

// Instance _
type Instance struct {
	ctx       context.Context
	tmpPath   files.Pather
	dealer    models.IWorkDealer
	logger    logrus.FieldLogger
	closed    chan struct{}
	performer models.IAuthor
}

// NewWorker _
func NewWorker(ctx context.Context, tmpPath files.Pather, dealer models.IWorkDealer) (*Instance, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.worker"); logger == nil {
		logger = ctxlog.New("fftb.worker")
	}

	performer, err := dealer.AllocatePerformerAuthority(uuid.New().String())

	if err != nil {
		return nil, errors.Wrap(err, "Obtaining performer authority")
	}

	return &Instance{
		ctx:       ctx,
		tmpPath:   tmpPath,
		dealer:    dealer,
		logger:    logger,
		performer: performer,
		closed:    make(chan struct{}),
	}, nil
}

// Start _
func (w *Instance) Start() (done chan struct{}) {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				close(w.closed)
				return
			default:
				freeSegment, err := w.dealer.FindFreeSegment(w.performer)

				if err != nil {
					if errors.Is(err, models.ErrNotFound) {
						w.logger.Debug("Free segment not found")
						time.Sleep(FreeTaskDelay)
					} else {
						w.logger.WithError(err).Warn("Searching free segment error")
					}

					continue
				}

				slog := w.logger.WithField(dlog.KeySegmentID, freeSegment.GetID()).
					WithField(dlog.KeyOrderID, freeSegment.GetOrderID()).
					WithField(dlog.KeyPerformer, w.performer.GetName())

				slog.Info("Found free segment")

				err = w.proceedSegment(slog, freeSegment)

				if err != nil {
					slog.WithError(err).Warn("Processing segment error")
				}
			}
			// freeSegment, err := w.dealer.FindFreeSegment(w.performer)

		}
	}()

	return w.closed
}

// proceedSegment _
func (w *Instance) proceedSegment(slog *logrus.Entry, freeSegment models.ISegment) error {
	convertSegment := freeSegment.(*models.ConvertSegment)

	inputClaim, err := w.dealer.GetInputStorageClaim(w.performer, freeSegment.GetID())

	if err != nil {
		panic(errors.Wrap(err, "Getting input storage claim"))
	}

	outputClaim, err := w.dealer.AllocateOutputStorageClaim(w.performer, freeSegment.GetID())

	if err != nil {
		panic(errors.Wrap(err, "Getting output storage claim"))
	}

	inputFile := w.tmpPath.BuildFile(fmt.Sprintf("%s.%s", inputClaim.GetID(), convertSegment.Muxer))

	err = inputFile.Create()

	if err != nil {
		panic(errors.Wrap(err, "Getting output storage claim"))
	}

	inputFileWriter, err := inputFile.WriteContent()

	if err != nil {
		panic(errors.Wrap(err, "Building input file writer"))
	}

	inputClaimReader, err := inputClaim.GetReader()

	if err != nil {
		panic(errors.Wrap(err, "Building input claim reader"))
	}

	outputClaimWriter, err := outputClaim.GetWriter()

	if err != nil {
		panic(errors.Wrap(err, "Building output claim reader"))
	}

	inputSize, err := inputClaim.GetSize()

	if err != nil {
		panic(errors.Wrap(err, "Getting input size"))
	}

	iopInputClaimReader := progress.NewReader(inputClaimReader)
	iopProgressChan := progress.NewTicker(w.ctx, iopInputClaimReader, int64(inputSize), 1*time.Second)

	go func() {
		for p := range iopProgressChan {
			err = w.dealer.NotifyRawDownload(w.performer, freeSegment.GetID(), makeIoProgresser(p, models.DownloadingInputStep))

			if err != nil {
				slog.WithError(err).Warn("Problem with notifying raw download")
			}
		}
	}()

	_, err = io.Copy(inputFileWriter, iopInputClaimReader)

	if err != nil {
		panic(errors.Wrap(err, "Writing segment from storage claim"))
	}

	outputFile := w.tmpPath.BuildFile(fmt.Sprintf("%s.%s", outputClaim.GetID(), convertSegment.Muxer))

	task := convert.Task{
		InFile:  inputFile.FullPath(),
		OutFile: outputFile.FullPath(),
		Params:  convertSegment.Params,
	}

	infoGetter := info.New()
	converter := convert.NewConverter(w.ctx, infoGetter)

	throttled := throttle.New(2000 * time.Millisecond)

	progressChan, doneChan, errChan := converter.Convert(task)

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("Terminating worker thread")

			<-converter.Closed()

			err = w.dealer.QuitSegment(w.performer, freeSegment.GetID())

			if err != nil {
				slog.WithError(err).Warn("Problem with quiting segment")
			}

			err = inputFile.Remove()

			if err != nil {
				slog.WithError(err).Warn("Problem with removing input file")
			}

			err = outputFile.Remove()

			if err != nil {
				slog.WithError(err).Warn("Problem with removing output file")
			}

			return nil

		case pmsg, ok := <-progressChan:
			if ok {
				throttled(func() {
					modProgress := makeProgresserFromConvert(pmsg)

					err = w.dealer.NotifyProcess(w.performer, freeSegment.GetID(), modProgress)

					if err != nil {
						slog.WithError(err).Warn("Problem with notifying process")
					}

					dlog.SegmentProgress(w.logger, freeSegment, modProgress)
				})
			}

		case cErr, ok := <-errChan:
			if ok {
				panic(errors.Wrap(cErr, "Error processing convert task"))
			}

		case <-doneChan:
			outputFileReader, err := outputFile.ReadContent()

			if err != nil {
				panic(errors.Wrap(err, "Building output file reader"))
			}

			_, err = io.Copy(outputClaimWriter, outputFileReader)

			if err != nil {
				panic(errors.Wrap(err, "Writing result to output claim"))
			}

			err = outputFile.Remove()

			if err != nil {
				panic(errors.Wrap(err, "Removing output file after uploading to storage"))
			}

			err = inputFile.Remove()

			if err != nil {
				panic(errors.Wrap(err, "Removing input file after processing it"))
			}

			slog.Info("Segment is done")

			err = w.dealer.FinishSegment(w.performer, freeSegment.GetID())

			if err != nil {
				panic(errors.Wrap(err, "Sending segment finish notification"))
			}

			return nil
		}
	}
}

// Closed _
func (w *Instance) Closed() <-chan struct{} {
	return w.closed
}
