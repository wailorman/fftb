package media

import (
	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/pkg/files"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
	"github.com/wailorman/goffmpeg/transcoder"
)

// InfoGetterInstance _
type InfoGetterInstance struct {
}

// InfoGetter _
type InfoGetter interface {
	GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error)
}

// NewInfoGetter _
func NewInfoGetter() *InfoGetterInstance {
	return &InfoGetterInstance{}
}

// GetMediaInfo _
func (ig *InfoGetterInstance) GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error) {
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
