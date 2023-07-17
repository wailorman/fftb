package handlers

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/ffmpeg"
	"github.com/wailorman/fftb/pkg/throttle"
)

// ConvertTaskHandler _
type ConvertTaskHandler struct {
	ctx           context.Context
	cancel        func()
	task          *pb.Task
	dealer        pb.Dealer
	logger        *logrus.Entry
	inputPath     string
	outputPath    string
	throttled     throttle.Throttler
	authorization string
	dealerHelper  *dealerHelper
	rcloneHelper  *rcloneHelper

	workingDir   string
	ffmpegClient *ffmpeg.FFmpegClient
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
	logger := ctxlog.WithPrefix(params.Logger, "handlers/convert")

	h := &ConvertTaskHandler{
		ctx:           ctx,
		cancel:        cancel,
		dealer:        params.Dealer,
		logger:        logger,
		workingDir:    params.WorkingDir,
		task:          params.Task,
		inputPath:     filepath.Join(params.WorkingDir, "input"),
		outputPath:    filepath.Join(params.WorkingDir, "output"),
		throttled:     throttle.New(5 * time.Second),
		authorization: params.Authorization,
	}

	h.dealerHelper = newDealerHelper(h)

	h.rcloneHelper = newRcloneHelper(h, rcloneHelperParams{
		localRemotesMap:  params.LocalRemotesMap,
		rcloneConfigPath: params.RcloneConfigPath,
		rclonePath:       params.RclonePath,
	})

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

func (h *ConvertTaskHandler) Run() error {
	if h.task.Type != pb.Task_CONVERT_V1 {
		return errors.New("Unexpected h.task type: `" + h.task.Type.String() + "`")
	}

	if err := h.rcloneHelper.touch(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.rcloneHelper.pull(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.convert(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.rcloneHelper.push(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	h.dealerHelper.exit(nil)
	return nil
}

func (h *ConvertTaskHandler) convert() error {
	outputPaths := searchOutputDirs(h.workingDir, h.task.ConvertParams.Opts)

	for _, path := range outputPaths {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.Wrapf(err, "Creating output directory `%s`", path)
		}
	}

	progress := make(chan *pb.ConvertTaskProgress)

	go func() { // TODO: ffmpeg message timeout
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

func searchOutputDirs(workingDir string, opts []string) []string {
	r := regexp.MustCompile(`^output/.+/.+`)
	result := make([]string, 0)

	for _, opt := range opts {
		if r.MatchString(opt) {
			result = append(result, (filepath.Join(workingDir, filepath.Dir(opt))))
		}
	}

	return result
}
