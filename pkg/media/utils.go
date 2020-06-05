package media

import (
	"github.com/wailorman/ffchunker/pkg/files"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
)

func filterVideos(allFiles []*files.File, infoGetter InfoGetter) []files.Filer {
	// c.VideoFileFiltered = make(chan VideoFileFilteringMessage, len(allFiles))

	videoFiles := make([]files.Filer, 0)

	for _, file := range allFiles {
		mediaInfo, err := infoGetter.GetMediaInfo(file)

		if err != nil {
			continue
		}

		if isVideo(mediaInfo) {
			videoFiles = append(videoFiles, file)
		}
	}

	return videoFiles
}

func isVideo(metadata ffmpegModels.Metadata) bool {
	if len(metadata.Streams) == 0 {
		return false
	}

	if metadata.Streams[0].BitRate == "" {
		return false
	}

	return true
}

func getVideoCodec(metadata ffmpegModels.Metadata) string {
	if !isVideo(metadata) {
		return ""
	}

	return metadata.Streams[0].CodecName
}
