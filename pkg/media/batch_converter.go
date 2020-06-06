package media

import (
	"strconv"
	"sync"
)

// BatchConverter _
type BatchConverter struct {
	ConversionStarted       chan bool
	TaskConversionStarted   chan ConverterTask
	MetadataReceived        chan MetadataReceivedBatchMessage
	InputVideoCodecDetected chan InputVideoCodecDetectedBatchMessage
	ConversionStopping      chan ConverterTask
	ConversionStopped       chan ConverterTask
	VideoFileFiltered       chan BatchVideoFilteringMessage

	infoGetter     InfoGetter
	stopConversion chan struct{}
}

// NewBatchConverter _
func NewBatchConverter(infoGetter InfoGetter) *BatchConverter {
	return &BatchConverter{
		infoGetter:     infoGetter,
		stopConversion: make(chan struct{}),
	}
}

// StopConversion _
func (bc *BatchConverter) StopConversion() {
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
					err := bc.convertOne(task, progress)

					if err != nil {
						failed <- BatchErrorMessage{
							Task: task,
							Err:  err,
						}

						if batchTask.StopConversionOnError {
							bc.StopConversion()
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
	singleConverter := NewConverter(bc.infoGetter)
	_progress, _finished, _failed := singleConverter.Convert(task)

	for {
		select {
		case <-bc.stopConversion:
			singleConverter.StopConversion()

		case <-singleConverter.ConversionStarted:
			bc.TaskConversionStarted <- task

		case metadata := <-singleConverter.MetadataReceived:
			bc.MetadataReceived <- MetadataReceivedBatchMessage{
				Metadata: metadata,
				Task:     task,
			}

		case videoCodec := <-singleConverter.InputVideoCodecDetected:
			bc.InputVideoCodecDetected <- InputVideoCodecDetectedBatchMessage{
				Codec: videoCodec,
				Task:  task,
			}

		case <-singleConverter.ConversionStopping:
			bc.ConversionStopping <- task

		case <-singleConverter.ConversionStopped:
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
