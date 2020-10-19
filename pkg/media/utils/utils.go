package utils

import (
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"

	"github.com/wailorman/fftb/pkg/files"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
)

// FilterVideos _
func FilterVideos(allFiles []*files.File, infoGetter mediaInfo.Getter) []files.Filer {
	videoFiles := make([]files.Filer, 0)

	for _, file := range allFiles {
		mediaInfo, err := infoGetter.GetMediaInfo(file)

		if err != nil {
			continue
		}

		if IsVideo(mediaInfo) {
			videoFiles = append(videoFiles, file)
		}
	}

	return videoFiles
}

// IsVideo _
func IsVideo(metadata ffmpegModels.Metadata) bool {
	if len(metadata.Streams) == 0 {
		return false
	}

	if metadata.Streams[0].BitRate == "" {
		return false
	}

	return true
}

// GetVideoCodec _
func GetVideoCodec(metadata ffmpegModels.Metadata) string {
	if !IsVideo(metadata) {
		return ""
	}

	return metadata.Streams[0].CodecName
}
