package convert

import (
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	mediaConvert "github.com/wailorman/fftb/pkg/media/convert"
)

func logProgress(msg mediaConvert.BatchProgressMessage) {
	progress := msg.Progress

	ctxlog.Logger.WithFields(logrus.Fields{
		"id":               msg.Task.ID,
		"frames_processed": progress.FramesProcessed(),
		"current_time":     progress.CurrentTime(),
		"current_bitrate":  progress.CurrentBitrate(),
		"progress":         progress.Progress(),
		"speed":            progress.Speed(),
		"fps":              progress.FPS(),
		"file_path":        progress.File().FullPath(),
	}).Info("Converting progress")
}

func logError(errorMessage mediaConvert.BatchErrorMessage) {
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
