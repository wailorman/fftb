package convert

import (
	"strconv"
	"sync"

	"github.com/wailorman/fftb/pkg/media/info"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
)

// BatchConverter _
type BatchConverter struct {
	ConversionStarted       chan bool
	TaskConversionStarted   chan ConverterTask
	MetadataReceived        chan MetadataReceivedBatchMessage
	InputVideoCodecDetected chan InputVideoCodecDetectedBatchMessage
	ConversionStopping      chan ConverterTask
	ConversionStopped       chan ConverterTask

	infoGetter           info.Getter
	stopConversion       chan struct{}
	conversionWasStopped bool
}

// NewBatchConverter _
func NewBatchConverter(infoGetter mediaInfo.Getter) *BatchConverter {
	return &BatchConverter{
		infoGetter:     infoGetter,
		stopConversion: make(chan struct{}),
	}
}

// Stop _
func (bc *BatchConverter) Stop() {
	bc.stopConversion = make(chan struct{})
	bc.conversionWasStopped = true
	// broadcast to all channel receivers
	close(bc.stopConversion)
}

// initChannels _
func (bc *BatchConverter) initChannels(taskCount int) {
	bc.ConversionStarted = make(chan bool, 1)
	bc.TaskConversionStarted = make(chan ConverterTask, taskCount)
	bc.MetadataReceived = make(chan MetadataReceivedBatchMessage, taskCount)
	bc.InputVideoCodecDetected = make(chan InputVideoCodecDetectedBatchMessage, taskCount)
	bc.ConversionStopping = make(chan ConverterTask, taskCount)
	bc.ConversionStopped = make(chan ConverterTask, taskCount)
}

// closeChannels _
func (bc *BatchConverter) closeChannels() {
	close(bc.ConversionStarted)
	close(bc.TaskConversionStarted)
	close(bc.MetadataReceived)
	close(bc.InputVideoCodecDetected)
	close(bc.ConversionStopping)
	close(bc.ConversionStopped)
}

// Convert _
func (bc *BatchConverter) Convert(batchTask BatchConverterTask) (
	progress chan BatchProgressMessage,
	finished chan bool,
	failed chan BatchErrorMessage,
) {
	taskCount := len(batchTask.Tasks)

	progress = make(chan BatchProgressMessage)
	finished = make(chan bool)
	failed = make(chan BatchErrorMessage)
	taskQueue := make(chan ConverterTask, taskCount)
	bc.initChannels(taskCount)

	var wg sync.WaitGroup
	wg.Add(taskCount)

	go func() {
		bc.ConversionStarted <- true

		for i := 0; i < batchTask.Parallelism; i++ {
			go func() {
				for task := range taskQueue {
					if !bc.conversionWasStopped {
						err := bc.convertOne(task, progress)

						if err != nil {
							failed <- BatchErrorMessage{
								Task: task,
								Err:  err,
							}

							if batchTask.StopConversionOnError {
								bc.Stop()
							}
						}
					}

					wg.Done()
				}
			}()
		}
	}()

	go func() {
		for i, task := range batchTask.Tasks {
			if task.ID == "" {
				task.ID = strconv.Itoa(i)
			}

			taskQueue <- task
		}
	}()

	go func() {
		wg.Wait()
		finished <- true

		close(progress)
		close(finished)
		close(failed)
		close(taskQueue)
		bc.closeChannels()
	}()

	return progress, finished, failed
}

func (bc *BatchConverter) convertOne(task ConverterTask, progress chan BatchProgressMessage) error {
	sConv := NewConverter(bc.infoGetter)
	_progress, _finished, _failed := sConv.Convert(task)

	for {
		select {
		case <-bc.stopConversion:
			sConv.Stop()

		case <-sConv.ConversionStarted:
			bc.TaskConversionStarted <- task

		case metadata := <-sConv.MetadataReceived:
			bc.MetadataReceived <- MetadataReceivedBatchMessage{
				Metadata: metadata,
				Task:     task,
			}

		case videoCodec := <-sConv.InputVideoCodecDetected:
			bc.InputVideoCodecDetected <- InputVideoCodecDetectedBatchMessage{
				Codec: videoCodec,
				Task:  task,
			}

		case <-sConv.ConversionStopping:
			bc.ConversionStopping <- task

		case <-sConv.ConversionStopped:
			bc.ConversionStopped <- task

		case progressMessage := <-_progress:
			progress <- BatchProgressMessage{
				Progress: progressMessage,
				Task:     task,
			}

		case errorMessage := <-_failed:
			return errorMessage

		case <-_finished:
			return nil

		}
	}
}
