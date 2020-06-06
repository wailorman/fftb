package convert

import (
	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker/pkg/media"
)

func logProgress(log *logrus.Entry, msg media.BatchProgressMessage) {
	progress := msg.Progress

	log.WithFields(logrus.Fields{
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

func logError(log *logrus.Entry, errorMessage media.BatchErrorMessage) {
	if errorMessage.Err != nil {
		log.WithField("error", errorMessage.Err.Error()).
			WithField("task_id", errorMessage.Task.ID).
			WithField("task_input_file", errorMessage.Task.InFile.FullPath()).
			Warn("Error")
	}
}

func logDone(log *logrus.Entry) {
	log.Info("Conversion done")
}

func logConversionStarted(log *logrus.Entry) {
	log.Info("Conversion started")
}

func logInputVideoCodec(log *logrus.Entry, msg media.InputVideoCodecDetectedBatchMessage) {
	log.WithField("input_video_codec", msg.Codec).
		WithField("task_id", msg.Task.ID).
		WithField("task_input_file", msg.Task.InFile.FullPath()).
		Debug("Input video codec detected")
}

func logTaskConversionStarted(log *logrus.Entry, task media.ConverterTask) {
	log.WithField("task_id", task.ID).
		WithField("task_input_file", task.InFile.FullPath()).
		Debug("Task conversion started")
}
