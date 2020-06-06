package convert

import (
	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/media"
)

func logProgress(msg media.BatchProgressMessage) {
	progress := msg.Progress

	ctxlog.Logger.WithFields(logrus.Fields{
		"id":               msg.Task.ID,
		"frames_processed": progress.FramesProcessed,
		"current_time":     progress.CurrentTime,
		"current_bitrate":  progress.CurrentBitrate,
		"progress":         progress.Progress,
		"speed":            progress.Speed,
		"fps":              progress.FPS,
		"file_path":        progress.File.FullPath(),
	}).Info("Converting progress")
}

func logError(errorMessage media.BatchErrorMessage) {
	if errorMessage.Err != nil {
		ctxlog.Logger.WithField("error", errorMessage.Err.Error()).
			WithField("task_id", errorMessage.Task.ID).
			WithField("task_input_file", errorMessage.Task.InFile.FullPath()).
			Warn("Error")
	}
}

func logDone() {
	ctxlog.Logger.Info("Conversion done")
}

func logConversionStarted() {
	ctxlog.Logger.Info("Conversion started")
}

func logInputVideoCodec(msg media.InputVideoCodecDetectedBatchMessage) {
	ctxlog.Logger.WithField("input_video_codec", msg.Codec).
		WithField("task_id", msg.Task.ID).
		WithField("task_input_file", msg.Task.InFile.FullPath()).
		Debug("Input video codec detected")
}

func logTaskConversionStarted(task media.ConverterTask) {
	ctxlog.Logger.WithField("task_id", task.ID).
		WithField("task_input_file", task.InFile.FullPath()).
		Debug("Task conversion started")
}
