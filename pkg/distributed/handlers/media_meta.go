package handlers

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/run"
	"github.com/wailorman/fftb/pkg/throttle"
)

type MediaMetaHandler struct {
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

	ffprobePath string
	workingDir  string
}

type MediaMetaHandlerParams struct {
	Ctx           context.Context
	Logger        *logrus.Entry
	WorkingDir    string
	Task          *pb.Task
	Dealer        pb.Dealer
	Authorization string

	RcloneConfigPath string
	RclonePath       string
	FFprobePath      string
	LocalRemotesMap  map[string]string
}

func NewMediaMetaHandler(params MediaMetaHandlerParams) *MediaMetaHandler {
	ctx, cancel := context.WithCancel(params.Ctx)
	logger := ctxlog.WithPrefix(params.Logger, "handlers/media_meta")

	h := &MediaMetaHandler{
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

	if params.FFprobePath != "" {
		h.ffprobePath = params.FFprobePath
	} else {
		h.ffprobePath = "ffprobe"
	}

	h.dealerHelper = newDealerHelper(h)

	h.rcloneHelper = newRcloneHelper(h, rcloneHelperParams{
		localRemotesMap:  params.LocalRemotesMap,
		rcloneConfigPath: params.RcloneConfigPath,
		rclonePath:       params.RclonePath,
	})

	return h
}

func (h *MediaMetaHandler) Run() error {
	if h.task.Type != pb.Task_MEDIA_META_V1 {
		return errors.New("Unexpected task type: `" + h.task.Type.String() + "`")
	}

	if err := h.rcloneHelper.touch(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.rcloneHelper.pull(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.extract(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	if err := h.rcloneHelper.push(); err != nil {
		h.dealerHelper.exit(err)
		return nil
	}

	h.logger.Trace("Task completed")

	h.dealerHelper.exit(nil)
	return nil
}

func (h *MediaMetaHandler) extract() error {
	fileBaseName := filepath.Base(h.task.MediaMetaParams.InputRclonePath)
	fileBaseNameNoExt := strings.TrimSuffix(fileBaseName, filepath.Ext(fileBaseName))

	operation := run.New([]string{
		h.ffprobePath,
		"-v", "quiet", "-print_format", "json", "-show_format", "-show_streams",
		filepath.Join(h.inputPath, fileBaseName),
	})
	operation.SetLogger(h.logger)

	if err := operation.Run(h.ctx); err != nil {
		return errors.Wrap(err, "Enqueueing ffprobe run")
	}

	stdout, _, err := operation.WaitOutput()

	if err != nil {
		return errors.Wrap(err, "Running ffprobe")
	}

	outputPath := filepath.Join(h.outputPath, fileBaseNameNoExt+".json")

	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		errors.Wrap(err, "Creating output file")
	}

	dataWriter := bufio.NewWriter(file)

	for _, data := range stdout {
		dataWriter.WriteString(data + "\n")
	}

	dataWriter.Flush()
	file.Close()

	return nil
}
