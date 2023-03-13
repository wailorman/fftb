package handlers

import (
	"context"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/ffmpeg"
	"github.com/wailorman/fftb/pkg/rclone"
	"github.com/wailorman/fftb/pkg/throttle"
)

// ConvertTaskHandler _
type ConvertTaskHandler struct {
	ctx           context.Context
	cancel        func()
	task          *pb.Task
	dealer        pb.Dealer
	logger        *logrus.Entry
	workingDir    string
	rclone        *rclone.RcloneClient
	ffmpegClient  *ffmpeg.FFmpegClient
	inputPath     string
	outputPath    string
	throttled     throttle.Throttler
	authorization string
}

type ConvertTaskHandlerParams struct {
	Ctx           context.Context
	Logger        *logrus.Entry
	WorkingDir    string
	Task          *pb.Task
	Dealer        pb.Dealer
	Authorization string

	RcloneConfigPath string
	RclonePath       string
	FFmpegPath       string
	FFprobePath      string
	LocalRemotesMap  map[string]string
}

func NewConvertTaskHandler(params ConvertTaskHandlerParams) *ConvertTaskHandler {
	ctx, cancel := context.WithCancel(params.Ctx)

	h := &ConvertTaskHandler{
		ctx:           ctx,
		cancel:        cancel,
		task:          params.Task,
		dealer:        params.Dealer,
		logger:        ctxlog.WithPrefix(params.Logger, "handlers/convert"),
		workingDir:    params.WorkingDir,
		inputPath:     filepath.Join(params.WorkingDir, "input"),
		outputPath:    filepath.Join(params.WorkingDir, "output"),
		throttled:     throttle.New(5 * time.Second),
		authorization: params.Authorization,
	}

	h.rclone = rclone.NewRcloneClient(rclone.RcloneClientParams{
		LocalRemotesMap: params.LocalRemotesMap,
	})
	h.rclone.SetLogger(params.Logger)

	if params.RcloneConfigPath != "" {
		h.rclone.SetConfigPath(params.RcloneConfigPath)
	}

	if params.RclonePath != "" {
		h.rclone.SetPath(params.RclonePath)
	}

	h.ffmpegClient = ffmpeg.NewFFmpegClient()
	h.ffmpegClient.SetLogger(h.logger)
	h.ffmpegClient.SetWorkingDir(h.workingDir)

	if params.FFmpegPath != "" {
		h.ffmpegClient.SetFFmpegPath(params.FFmpegPath)
	}

	if params.FFprobePath != "" {
		h.ffmpegClient.SetFFprobePath(params.FFprobePath)
	}

	return h
}

// Run _
func (h *ConvertTaskHandler) Run() error {
	if h.task.Type != pb.TaskType_CONVERT_V1 {
		return errors.New("Unexpected task type: `" + h.task.Type.String() + "`")
	}

	var err error
	var isInputLocal, isOutputLocal bool

	if isInputLocal, err = h.rclone.Touch(h.task.ConvertParams.InputRclonePath, h.inputPath); err != nil {
		h.exit(errors.Wrap(err, "Touching input path"))
		return nil
	}

	if isOutputLocal, err = h.rclone.Touch(h.task.ConvertParams.OutputRclonePath, h.outputPath); err != nil {
		h.exit(errors.Wrap(err, "Touching output path"))
		return nil
	}

	if !isInputLocal {
		if err := h.pull(); err != nil {
			h.exit(err)
			return nil
		}
	}

	if err := h.convert(); err != nil {
		h.exit(err)
		return nil
	}

	if !isOutputLocal {
		if err := h.push(); err != nil {
			h.exit(err)
			return nil
		}
	}

	h.exit(nil)
	return nil
}

func (h *ConvertTaskHandler) pull() error {
	progressCh := make(chan rclone.ProgressMessage)

	go func() {
		for progressMessage := range progressCh {
			if progressMessage.IsValid() {
				h.throttled(func() {
					_, err := h.dealer.Notify(h.ctx, &pb.NotifyRequest{
						Step:          pb.NotifyRequest_DOWNLOADING_INPUT,
						Authorization: h.authorization,
						TaskId:        h.task.Id,
						Progress:      progressMessage.Progress(),
					})

					if err != nil && !errors.Is(err, context.Canceled) {
						h.logger.
							WithError(err).
							Warn("Failed to notify pull")

						h.cancel()
						return
					}

					h.logger.
						WithField(dlog.KeyProgress, progressMessage.Progress()).
						WithField(dlog.KeySpeed, progressMessage.HumanSpeed()).
						Info("Downloading input")
				})
			}
		}
	}()

	err := h.rclone.Pull(h.ctx, h.task.ConvertParams.InputRclonePath, h.inputPath, progressCh)

	if err != nil && !errors.Is(err, context.Canceled) {
		h.logger.WithError(err).Warn("Failed to pull input")
	}

	return err
}

func (h *ConvertTaskHandler) convert() error {
	progress := make(chan *pb.ConvertTaskProgress)

	go func() {
		for progressMessage := range progress {
			h.throttled(func() {
				_, err := h.dealer.Notify(h.ctx, &pb.NotifyRequest{
					Step:            pb.NotifyRequest_PROCESSING,
					Authorization:   h.authorization,
					TaskId:          h.task.Id,
					Progress:        0,
					ConvertProgress: progressMessage,
				})

				if err != nil && !errors.Is(err, context.Canceled) {
					h.logger.
						WithError(err).
						Warn("Failed to notify convert")

					h.cancel()
					return
				}

				h.logger.
					WithField(dlog.KeyFPS, progressMessage.Fps).
					Info("Converting")
			})
		}
	}()

	err := h.ffmpegClient.Transcode(h.ctx, h.task.ConvertParams.Opts, progress)

	if err != nil && !errors.Is(err, context.Canceled) {
		h.logger.WithError(err).Warn("Failed to convert")
	}

	return err
}

func (h *ConvertTaskHandler) push() error {
	progressCh := make(chan rclone.ProgressMessage)

	go func() {
		for progressMessage := range progressCh {
			if progressMessage.IsValid() {
				h.throttled(func() {
					_, err := h.dealer.Notify(h.ctx, &pb.NotifyRequest{
						Step:          pb.NotifyRequest_UPLOADING_OUTPUT,
						Authorization: h.authorization,
						TaskId:        h.task.Id,
						Progress:      progressMessage.Progress(),
					})

					if err != nil && !errors.Is(err, context.Canceled) {
						h.logger.
							WithError(err).
							Warn("Failed to notify push")

						h.cancel()
						return
					}

					h.logger.
						WithField(dlog.KeyProgress, progressMessage.Progress()).
						WithField(dlog.KeySpeed, progressMessage.HumanSpeed()).
						Info("Uploading output")
				})
			}
		}
	}()

	err := h.rclone.Push(h.ctx, h.outputPath, h.task.ConvertParams.OutputRclonePath, progressCh)

	if err != nil && !errors.Is(err, context.Canceled) {
		h.logger.WithError(err).Warn("Failed to push output")
	}

	return err
}

func (h *ConvertTaskHandler) exit(err error) {
	if errors.Is(h.ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		h.quit()
		return
	}

	if err != nil {
		h.fail(err)
		return
	}

	h.finish()
}

func (h *ConvertTaskHandler) finish() {
	_, err := h.dealer.FinishTask(context.TODO(), &pb.FinishTaskRequest{
		Authorization: h.authorization,
		TaskId:        h.task.Id,
	})

	if err != nil {
		h.logger.WithError(err).Warn("Failed to finish")
	}
}

func (h *ConvertTaskHandler) quit() {
	_, err := h.dealer.QuitTask(context.TODO(), &pb.QuitTaskRequest{
		Authorization: h.authorization,
		TaskId:        h.task.Id,
	})

	if err != nil {
		h.logger.WithError(err).Warn("Failed to quit")
	}
}

func (h *ConvertTaskHandler) fail(failure error) {
	_, err := h.dealer.FailTask(h.ctx, failTaskRequest(h.authorization, h.task.Id, failure))

	if err != nil {
		h.logger.WithError(err).Warn("Failed to report failure")
	}
}

func failTaskRequest(authorization string, taskId string, err error) *pb.FailTaskRequest {
	return &pb.FailTaskRequest{
		Authorization: authorization,
		TaskId:        taskId,
		Failures:      []string{err.Error()},
	}
}
