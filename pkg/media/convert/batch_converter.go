package convert

import (
	"context"
	"strconv"

	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/media/info"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
)

// BatchConverter _
type BatchConverter struct {
	wg         *chwg.ChannelledWaitGroup
	ctx        context.Context
	infoGetter info.Getter
}

// NewBatchConverter _
func NewBatchConverter(ctx context.Context, infoGetter mediaInfo.Getter) *BatchConverter {
	return &BatchConverter{
		wg:         chwg.New(),
		ctx:        ctx,
		infoGetter: infoGetter,
	}
}

// Convert _
func (bc *BatchConverter) Convert(batchTask BatchConverterTask) (
	progress chan BatchProgressMessage,
	failures chan BatchErrorMessage,
) {
	taskCount := len(batchTask.Tasks)

	progress = make(chan BatchProgressMessage)
	failures = make(chan BatchErrorMessage)
	taskQueue := make(chan ConverterTask, taskCount)

	bc.wg.Add(taskCount)

	go func() {
		for i := 0; i < batchTask.Parallelism; i++ {
			go func() {
				for task := range taskQueue {
					if !bc.wg.IsFinished() {
						err := bc.convertOne(task, progress)

						if err != nil {
							failures <- BatchErrorMessage{
								Task: task,
								Err:  err,
							}

							if batchTask.StopConversionOnError {
								bc.wg.AllDone()
								return
							}
						}

						bc.wg.Done()
					}
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
		bc.wg.Wait()
		close(progress)
		close(failures)
		close(taskQueue)
	}()

	return progress, failures
}

func (bc *BatchConverter) convertOne(task ConverterTask, progress chan BatchProgressMessage) error {
	sCtx, sCancel := context.WithCancel(bc.ctx)
	sConv := NewConverter(sCtx, bc.infoGetter)
	sProgress, sFailures := sConv.Convert(task)

	for {
		select {
		case progressMessage, ok := <-sProgress:
			if ok {
				progress <- BatchProgressMessage{
					Progress: progressMessage,
					Task:     task,
				}
			}

		case errorMessage, failed := <-sFailures:
			sCancel()

			if !failed {
				<-sConv.Closed()
				return nil
			}

			return errorMessage
		}
	}
}

// Closed returns channel with finished signal
func (bc *BatchConverter) Closed() <-chan struct{} {
	return bc.wg.Closed()
}
