package info

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// Instance _
type Instance struct {
}

// Getter _
type Getter interface {
	GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error)
}

// New _
func New() Getter {
	return &Instance{}
}

// GetMediaInfo _
func (ig *Instance) GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error) {
	trans := &transcoder.Transcoder{}

	err := trans.InitializeEmptyTranscoder()

	if err != nil {
		return ffmpegModels.Metadata{}, errors.Wrap(err, "Initializing ffprobe instance")
	}

	metadata, err := trans.GetFileMetadata(file.FullPath())

	if err != nil {
		return ffmpegModels.Metadata{}, errors.Wrap(err, "Getting file metadata from ffprobe")
	}

	return metadata, nil
}
