package ffmpeg

import (
	"bytes"
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/goffmpeg/utils"
)

// LoggingPrefix _
const LoggingPrefix = "goffmpeg"

// Configuration ...
type Configuration struct {
	FfmpegBin  string
	FfprobeBin string
}

// Configure Get and set FFmpeg and FFprobe bin paths
func Configure(ctx context.Context) (Configuration, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, LoggingPrefix); logger == nil {
		logger = ctxlog.New(LoggingPrefix)
	}

	var outFFmpeg bytes.Buffer
	var outFFprobe bytes.Buffer

	execFFmpegCommand := utils.GetFFmpegExec()
	execFFprobeCommand := utils.GetFFprobeExec()

	outFFmpeg, err := utils.TestCmd(execFFmpegCommand[0], execFFmpegCommand[1])
	if err != nil {
		return Configuration{}, err
	}

	outFFprobe, err = utils.TestCmd(execFFprobeCommand[0], execFFprobeCommand[1])
	if err != nil {
		return Configuration{}, err
	}

	ffmpeg := strings.ReplaceAll(
		outFFmpeg.String(),
		utils.LineSeparator(),
		"",
	)

	ffprobe := strings.ReplaceAll(
		outFFprobe.String(),
		utils.LineSeparator(),
		"",
	)

	logger.WithFields(logrus.Fields{
		"ffmpeg_path":  ffmpeg,
		"ffprobe_path": ffprobe,
	}).Debug("Found ffmpeg binaries")

	cnf := Configuration{ffmpeg, ffprobe}
	return cnf, nil
}
