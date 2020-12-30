package split

import (
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	mediaConvert "github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/ff"
)

func logProgress(progress ff.Progressable) {
	ctxlog.Logger.WithFields(logrus.Fields{
		"frames_processed": progress.FramesProcessed(),
		"current_time":     progress.CurrentTime(),
		"current_bitrate":  progress.CurrentBitrate(),
		"progress":         progress.Progress(),
		"speed":            progress.Speed(),
		"fps":              progress.FPS(),
		"file_path":        progress.File().FullPath(),
	}).Info("Converting progress")
}

func logError(err error) {
	ctxlog.Logger.WithField("error", err.Error()).
		Warn("Error")
}

func logDone() {
	ctxlog.Logger.Info("Splitting done")
}

func logSplittingStarted() {
	ctxlog.Logger.Info("Splitting started")
}

func logTaskSplittingStarted(task mediaConvert.Task) {
	ctxlog.Logger.WithField("task_id", task.ID).
		WithField("task_input_file", task.InFile).
		Debug("Task splitting started")
}
