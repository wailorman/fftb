package worker

import (
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/info"
)

// FreeTaskDelay _
const FreeTaskDelay = time.Duration(10) * time.Second

// FreeTaskTimeout _
const FreeTaskTimeout = time.Duration(10) * time.Second

// Instance _
type Instance struct {
	tmpPath files.Pather
	dealer  models.IWorkDealer
}

// NewWorker _
func NewWorker(tmpPath files.Pather, dealer models.IWorkDealer) *Instance {
	return &Instance{
		tmpPath: tmpPath,
		dealer:  dealer,
	}
}

// Start _
func (w *Instance) Start() error {
	// done := make(chan struct{}, 0)
	// progress := make(chan models.Progresser, 0)
	// failures := make(chan error, 0)

	// panic(models.ErrNotImplemented)
	// TODO: stop polling delaer on cancel
	// freeSegment, err := w.dealer.FindFreeSegment()

	// go func() {
	for {
		freeSegment, err := w.dealer.FindFreeSegment("local")

		convertSegment := freeSegment.(*models.ConvertSegment)

		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				log.Println("Free task not found")
				time.Sleep(FreeTaskDelay)
			} else {
				log.Printf("Searching free task error: %s\n", err)
			}
		}

		inputClaim, err := w.dealer.GetInputStorageClaim(freeSegment)

		if err != nil {
			// errF := errors.Wrap(err, "Getting output storage claim")
			panic(errors.Wrap(err, "Getting input storage claim"))
		}

		outputClaim, err := w.dealer.GetOutputStorageClaim(freeSegment)

		if err != nil {
			// errF := errors.Wrap(err, "Getting output storage claim")
			panic(errors.Wrap(err, "Getting output storage claim"))
		}

		inputFile := w.tmpPath.BuildFile(inputClaim.GetID())

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

		_, err = io.Copy(inputFileWriter, inputClaimReader)

		if err != nil {
			panic(errors.Wrap(err, "Writing segment from storage claim"))
		}

		outputFile := w.tmpPath.BuildFile(outputClaim.GetID())

		batchTask := convert.BatchConverterTask{
			Parallelism: 1,
			Tasks: []convert.ConverterTask{
				convert.ConverterTask{
					InFile:  inputFile,
					OutFile: outputFile,
					// HWAccel:      c.String("hwa"),
					VideoCodec: convertSegment.Params.VideoCodec,
					// Preset:       c.String("preset"),
					// VideoBitRate: c.String("video-bitrate"),
					VideoQuality: convertSegment.Params.VideoQuality,
					// Scale:        c.String("scale"),
				},
			},
		}

		infoGetter := info.New()
		converter := convert.NewBatchConverter(infoGetter)
		_, doneChan, errChan := converter.Convert(batchTask)

		for {
			select {
			case <-doneChan:
				outputFileReader, err := outputFile.ReadContent()

				if err != nil {
					panic(errors.Wrap(err, "Building output file reader"))
				}

				_, err = io.Copy(outputClaimWriter, outputFileReader)

				if err != nil {
					panic(errors.Wrap(err, "Writing result to output claim"))
				}

				log.Printf("segment %s is done!", freeSegment.GetID())
				return nil
			case bErr := <-errChan:
				panic(errors.Wrap(bErr.Err, "Error processing convert task"))
			}
		}

	}
	// }()

	// return nil
}

// Cancel _
func (w *Instance) Cancel() {
	panic(models.ErrNotImplemented)
}

func (w *Instance) findFreeTaskWithTimeout(timeout time.Duration) (models.ISegment, error) {
	freeSegment := make(chan models.ISegment)
	failures := make(chan error)

	go func() {
		fSeg, err := w.dealer.FindFreeSegment("local")

		if err != nil {
			failures <- err
			return
		}

		freeSegment <- fSeg
	}()

	// select {
	// case fSeg := <-freeSegment:
	// case failure := <-failures:
	// case <-time.After(FreeTaskTimeout):
	// }

	return nil, nil // TODO:
}
