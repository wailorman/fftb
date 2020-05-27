package ffchunker

import (
	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/files"

	"github.com/wailorman/goffmpeg/transcoder"
)

// DurationCalculator _
type DurationCalculator struct {
}

// VideoDurationCalculator _
type VideoDurationCalculator interface {
	Calculate(file files.Filer) (float64, error)
}

// NewDurationCalculator _
func NewDurationCalculator() *DurationCalculator {
	return &DurationCalculator{}
}

// Calculate _
func (d *DurationCalculator) Calculate(file files.Filer) (float64, error) {
	trans := &transcoder.Transcoder{}

	err := trans.InitializeEmptyTranscoder()

	if err != nil {
		return 0, errors.Wrap(err, "Initializing ffprobe instance")
	}

	metadata, err := trans.GetFileMetadata(file.FullPath())

	if err != nil {
		return 0, errors.Wrap(err, "Getting file metadata from ffprobe")
	}

	if len(metadata.Streams) == 0 {
		return 0, errors.New("No streams in media file")
	}

	return metadata.Streams[0].DurationFloat, nil
}
