package worker

import (
	"context"
	"time"

	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"

	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/throttle"
	"github.com/wailorman/fftb/pkg/tmpmgr"
)

var throttledTaskNotFound = throttle.New(1 * time.Minute)

const freeTaskDelay = time.Duration(3) * time.Second

// Instance _
type Instance struct {
	ctx     context.Context
	tmpPath string
	logger  logrus.FieldLogger
	dealer  pb.Dealer
	wg      *chwg.ChannelledWaitGroup

	rcloneConfigPath string
	rclonePath       string
	ffmpegPath       string
	ffprobePath      string
}

type WorkerParams struct {
	Ctx    context.Context
	Dealer pb.Dealer
	Logger logrus.FieldLogger
	Wg     *chwg.ChannelledWaitGroup

	TmpPath          string
	RcloneConfigPath string
	RclonePath       string
	FFmpegPath       string
	FFprobePath      string
}

// NewWorker _
func NewWorker(params WorkerParams) *Instance {
	return &Instance{
		ctx:     params.Ctx,
		tmpPath: params.TmpPath,
		logger:  params.Logger,
		dealer:  params.Dealer,
		wg:      params.Wg,

		rcloneConfigPath: params.RcloneConfigPath,
		rclonePath:       params.RclonePath,
		ffmpegPath:       params.FFmpegPath,
		ffprobePath:      params.FFprobePath,
	}
}

// Start _
func (w *Instance) Start() {
	w.wg.Add(1)

	w.logger.Info("Worker started")

	tmpMgr := tmpmgr.New(w.tmpPath)

	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-w.ctx.Done():
				w.logger.Info("terminated")
				return

			default:
				freeTask, err := w.dealer.FindFreeTask(w.ctx, &pb.FindFreeTaskRequest{
					Authorization: "TODO:",
				})

				if err != nil {
					if twerr, ok := err.(twirp.Error); ok && twerr.Code() == twirp.NotFound {
						throttledTaskNotFound(func() {
							w.logger.Info("Free task not found")
						})
					} else {
						w.logger.WithError(err).Warn("Searching free task error")
					}

					time.Sleep(freeTaskDelay)
					continue
				}

				logger := w.logger.
					WithField(dlog.KeyTaskID, freeTask.Id).
					WithField(dlog.KeyRunID, uniuri.New())

				logger.Info("Found free task")

				tmpPath, err := tmpMgr.Create(freeTask.Id)

				if err != nil {
					logger.WithError(err).
						Fatal("Failed to create temporary directory")
				} else {
					logger.WithField(dlog.KeyPath, tmpPath).
						Debug("Created temporary directory")
				}

				convertHandler := handlers.NewConvertTaskHandler(handlers.ConvertTaskHandlerParams{
					Ctx:        w.ctx,
					Logger:     logger,
					WorkingDir: tmpPath,
					Task:       freeTask,
					Dealer:     w.dealer,

					RcloneConfigPath: w.rcloneConfigPath,
					RclonePath:       w.rclonePath,
					FFmpegPath:       w.ffmpegPath,
					FFprobePath:      w.ffprobePath,
				})

				w.wg.Add(1)

				if err = convertHandler.Run(); err != nil {
					logger.WithError(err).
						Warn("Failed to run convert handler")
				}

				if err = tmpMgr.Destroy(freeTask.Id); err != nil {
					logger.WithError(err).
						Warn("Failed to destroy temporary directory")

					return
				}

				w.wg.Done()
			}
		}
	}()
}
