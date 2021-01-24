package convert

import (
	"context"
	"strconv"
	"sync"

	"github.com/wailorman/fftb/pkg/media/info"
)

// BatchConverter _
type BatchConverter struct {
	closed     chan struct{}
	ctx        context.Context
	cancel     func()
	infoGetter info.Getter
}

// NewBatchConverter _
func NewBatchConverter(ctx context.Context, infoGetter info.Getter) *BatchConverter {
	bc := &BatchConverter{
		infoGetter: infoGetter,
		closed:     make(chan struct{}),
	}

	bc.ctx, bc.cancel = context.WithCancel(ctx)

	return bc
}

// Stop _
func (bc *BatchConverter) Stop() {
	bc.cancel()
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

	var wg sync.WaitGroup
	wg.Add(taskCount)

	go func() {
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

		close(bc.closed)

		select {
		case <-bc.ctx.Done():
		default:
			finished <- true
		}

		close(progress)
		close(finished)
		close(failed)
		close(taskQueue)
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

// Closed _
func (bc *BatchConverter) Closed() <-chan struct{} {
	return bc.closed
}
