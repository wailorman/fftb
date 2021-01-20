package convert

import (
	"context"
	"strconv"
	"sync"

	"github.com/wailorman/fftb/pkg/media/info"
)

// BatchConverter _
type BatchConverter struct {
	ConversionStarted       chan bool
	TaskConversionStarted   chan Task
	MetadataReceived        chan MetadataReceivedBatchMessage
	InputVideoCodecDetected chan InputVideoCodecDetectedBatchMessage
	ConversionStopping      chan Task
	ConversionStopped       chan Task

	ctx        context.Context
	cancel     func()
	infoGetter info.Getter
}

// NewBatchConverter _
func NewBatchConverter(infoGetter info.Getter) *BatchConverter {
	bc := &BatchConverter{
		infoGetter: infoGetter,
	}

	bc.ctx, bc.cancel = context.WithCancel(context.TODO())

	return bc
}

// Stop _
func (bc *BatchConverter) Stop() {
	bc.cancel()
}

// initChannels _
func (bc *BatchConverter) initChannels(taskCount int) {
	bc.ConversionStarted = make(chan bool, 1)
	bc.TaskConversionStarted = make(chan Task, taskCount)
	bc.MetadataReceived = make(chan MetadataReceivedBatchMessage, taskCount)
	bc.InputVideoCodecDetected = make(chan InputVideoCodecDetectedBatchMessage, taskCount)
	bc.ConversionStopping = make(chan Task, taskCount)
	bc.ConversionStopped = make(chan Task, taskCount)
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
func (bc *BatchConverter) Convert(batchTask BatchTask) (
	progress chan BatchProgressMessage,
	finished chan bool,
	failed chan BatchErrorMessage,
) {
	taskCount := len(batchTask.Tasks)

	progress = make(chan BatchProgressMessage)
	finished = make(chan bool)
	failed = make(chan BatchErrorMessage)
	taskQueue := make(chan Task, taskCount)
	bc.initChannels(taskCount)

	var wg sync.WaitGroup
	wg.Add(taskCount)

	go func() {
		bc.ConversionStarted <- true

		for i := 0; i < batchTask.Parallelism; i++ {
			go func() {
				for task := range taskQueue {
					select {
					case <-bc.ctx.Done():
						wg.Done()
						return
					default:
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

func (bc *BatchConverter) convertOne(task Task, progress chan BatchProgressMessage) error {
	sConv := NewConverter(bc.infoGetter)
	_progress, _finished, _failed := sConv.Convert(task)

	for {
		select {
		case <-bc.ctx.Done():
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
