package handlers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/rclone"
	"github.com/wailorman/fftb/pkg/throttle"
)

type rcloneHelper struct {
	dealer           pb.Dealer
	ctx              context.Context
	authorization    string
	logger           *logrus.Entry
	rclone           *rclone.RcloneClient
	throttled        throttle.Throttler
	cancel           func()
	inputLocalPath   string
	inputRemotePath  string
	outputLocalPath  string
	outputRemotePath string
	task             *pb.Task

	isInputLocal  bool
	isOutputLocal bool
}

type rcloneHelperParams struct {
	localRemotesMap  map[string]string
	rcloneConfigPath string
	rclonePath       string
}

func newRcloneHelper(handler interface{}, params rcloneHelperParams) *rcloneHelper {
	h := &rcloneHelper{}

	switch th := handler.(type) {
	case *ConvertTaskHandler:
		h.dealer = th.dealer
		h.ctx = th.ctx
		h.authorization = th.authorization
		h.logger = th.logger
		h.task = th.task
		h.throttled = th.throttled
		h.cancel = th.cancel

		h.inputLocalPath = th.inputPath
		h.inputRemotePath = th.task.ConvertParams.InputRclonePath
		h.outputLocalPath = th.outputPath
		h.outputRemotePath = th.task.ConvertParams.OutputRclonePath

	case *MediaMetaHandler:
		h.dealer = th.dealer
		h.ctx = th.ctx
		h.authorization = th.authorization
		h.logger = th.logger
		h.task = th.task
		h.throttled = th.throttled
		h.cancel = th.cancel

		h.inputLocalPath = th.inputPath
		h.inputRemotePath = th.task.MediaMetaParams.InputRclonePath
		h.outputLocalPath = th.outputPath
		h.outputRemotePath = th.task.MediaMetaParams.OutputRclonePath

	default:
		panic("Unknown handler type")
	}

	h.rclone = rclone.NewRcloneClient(rclone.RcloneClientParams{
		LocalRemotesMap: params.localRemotesMap,
	})

	h.rclone.SetLogger(h.logger)

	if params.rcloneConfigPath != "" {
		h.rclone.SetConfigPath(params.rcloneConfigPath)
	}

	if params.rclonePath != "" {
		h.rclone.SetPath(params.rclonePath)
	}

	return h
}

func (h *rcloneHelper) touch() error {
	var err error

	if h.isInputLocal, err = h.rclone.Touch(h.inputRemotePath, h.inputLocalPath); err != nil {
		return errors.Wrap(err, "Touching input path")
	}

	if h.isOutputLocal, err = h.rclone.Touch(h.outputRemotePath, h.outputLocalPath); err != nil {
		return errors.Wrap(err, "Touching output path")
	}

	return nil
}

func (h *rcloneHelper) pull() error {
	if h.isInputLocal {
		return nil
	}

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

	err := h.rclone.Pull(h.ctx, h.inputRemotePath, h.inputLocalPath, progressCh)

	if err != nil && !errors.Is(err, context.Canceled) {
		h.logger.WithError(err).Warn("Failed to pull input")
	}

	return err
}

func (h *rcloneHelper) push() error {
	if h.isOutputLocal {
		return nil
	}

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

	err := h.rclone.Push(h.ctx, h.outputLocalPath, h.outputRemotePath, progressCh)

	if err != nil && !errors.Is(err, context.Canceled) {
		h.logger.WithError(err).Warn("Failed to push output")
	}

	return err
}
