package minfo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// GetterInstance _
type GetterInstance struct {
}

// Getter _
type Getter interface {
	GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error)
}

// NewGetter _
func NewGetter() *GetterInstance {
	return &GetterInstance{}
}

// GetMediaInfo _
func (ig *GetterInstance) GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error) {
	trans := transcoder.New(context.TODO())

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
