package worker

import (
	"context"
	"fmt"
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
const FreeTaskDelay = time.Duration(3) * time.Second

// FreeTaskTimeout _
const FreeTaskTimeout = time.Duration(3) * time.Second

// Instance _
type Instance struct {
	ctx     context.Context
	tmpPath files.Pather
	dealer  models.IWorkDealer
	closed  chan struct{}
}

// NewWorker _
func NewWorker(ctx context.Context, tmpPath files.Pather, dealer models.IWorkDealer) *Instance {
	return &Instance{
		ctx:     ctx,
		tmpPath: tmpPath,
		dealer:  dealer,
		closed:  make(chan struct{}),
	}
}

// Start _
func (w *Instance) Start() error {
	// go func() {
	for {
		freeSegment, err := w.dealer.FindFreeSegment("local")

		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				log.Println("Free task not found")
				time.Sleep(FreeTaskDelay)
			} else {
				log.Printf("Searching free task error: %s\n", err)
			}

			continue
		}

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

	_, err = io.Copy(inputFileWriter, inputClaimReader)

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
	progressChan, doneChan, errChan := converter.Convert(batchTask)

	for {
		select {
		case pmsg := <-progressChan:
			fmt.Printf("pmsg: %#v\n", pmsg.Progress)

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
			return nil
		}
	}
}

// Closed _
func (w *Instance) Closed() <-chan struct{} {
	return w.closed
}
