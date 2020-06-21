package cut

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/goffmpeg/transcoder"
)

// CutterInstance _
type CutterInstance struct {
}

// NewCutter _
func NewCutter() *CutterInstance {
	return &CutterInstance{}
}

// Cutter _
type Cutter interface {
	CutVideo(inFile files.Filer, outFile files.Filer, offset float64, maxFileSize int) (files.Filer, error)
}

// CutVideo _
func (ci *CutterInstance) CutVideo(
	inFile files.Filer,
	outFile files.Filer,
	offset float64,
	maxFileSize int,
) (files.Filer, error) {

	log := ctxlog.Logger

	trans := new(transcoder.Transcoder)

	err := trans.Initialize(
		inFile.FullPath(),
		outFile.FullPath(),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Initializing ffmpeg transcoder")
	}

	trans.MediaFile().SetVideoCodec("copy")
	trans.MediaFile().SetAudioCodec("copy")
	trans.MediaFile().SetFileSizeLimit(strconv.Itoa(maxFileSize))
	trans.MediaFile().SetSeekTimeInput(fmt.Sprintf("%f", offset))

	done := trans.Run(true)

	progressChan := trans.Output()

	for {
		select {
		case progress := <-progressChan:
			log.WithFields(logrus.Fields{
				"frames_processed": progress.FramesProcessed,
				"current_time":     progress.CurrentTime,
				"current_bitrate":  progress.CurrentBitrate,
				"progress":         progress.Progress,
				"speed":            progress.Speed,
				"fps":              progress.FPS,
			}).Info("ffmpeg progress")

		case err := <-done:
			if err != nil {
				return nil, errors.Wrap(err, "ffmpeg cutting error")
			}

			return outFile, nil
		}
	}
}
