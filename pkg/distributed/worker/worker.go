package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	ctx     context.Context
	tmpPath files.Pather
	dealer  models.IWorkDealer
	logger  logrus.FieldLogger
	closed  chan struct{}
}

// NewWorker _
func NewWorker(ctx context.Context, tmpPath files.Pather, dealer models.IWorkDealer) *Instance {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.distributed.worker"); logger == nil {
		logger = ctxlog.New("fftb.distributed.worker")
	}

	return &Instance{
		ctx:     ctx,
		tmpPath: tmpPath,
		dealer:  dealer,
		logger:  logger,
		closed:  make(chan struct{}),
	}
}

// Start _
func (w *Instance) Start() error {
	// go func() {
	for {
		freeSegment, err := w.dealer.FindFreeSegment()

		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.logger.Debug("Free segment not found")
				time.Sleep(FreeTaskDelay)
			} else {
				w.logger.WithError(err).Warn("Searching free segment error")
			}

			continue
		}

		w.logger.
			WithField(dlog.KeySegmentID, freeSegment.GetID()).
			Info("Found free segment")

		err = w.proceedSegment(freeSegment)

		if err != nil {
			return errors.Wrap(err, "Processing segment")
		}
	}
}

// proceedSegment _
func (w *Instance) proceedSegment(freeSegment models.ISegment) error {
	convertSegment := freeSegment.(*models.ConvertSegment)

	inputClaim, err := w.dealer.GetInputStorageClaim(freeSegment)

	if err != nil {
		panic(errors.Wrap(err, "Getting input storage claim"))
	}

	outputClaim, err := w.dealer.AllocateOutputStorageClaim(freeSegment)

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
	iopProgressChan := progress.NewTicker(w.ctx, iopInputClaimReader, int64(inputSize), 100*time.Millisecond)

	go func() {
		for p := range iopProgressChan {
			w.dealer.NotifyRawDownload(freeSegment, makeIoProgresser(p, models.DownloadingInputStep))
		}
	}()

	_, err = io.Copy(inputFileWriter, iopInputClaimReader)

	if err != nil {
		panic(errors.Wrap(err, "Writing segment from storage claim"))
	}

	outputFile := w.tmpPath.BuildFile(fmt.Sprintf("%s.%s", outputClaim.GetID(), convertSegment.Muxer))

	batchTask := convert.BatchTask{
		Parallelism: 1,
		Tasks: []convert.Task{
			convert.Task{
				InFile:  inputFile.FullPath(),
				OutFile: outputFile.FullPath(),
				Params:  convertSegment.Params,
			},
		},
	}

	infoGetter := info.New()
	converter := convert.NewBatchConverter(infoGetter)

	throttled := throttle.New(2000 * time.Millisecond)

	progressChan, doneChan, errChan := converter.Convert(batchTask)

	for {
		select {
		case pmsg := <-progressChan:
			throttled(func() {
				modProgress := makeProgresserFromConvert(pmsg)

				w.dealer.NotifyProcess(freeSegment, modProgress)
				dlog.SegmentProgress(w.logger, freeSegment, modProgress)
			})

		case bErr := <-errChan:
			panic(errors.Wrap(bErr.Err, "Error processing convert task"))

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

			log.Printf("segment %s is done!", freeSegment.GetID())

			err = w.dealer.FinishSegment(freeSegment)

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
