package utils

import (
	"io"
	"os"

	"github.com/pkg/errors"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"

	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/info"
)

// FilterVideos _
func FilterVideos(allFiles []files.Filer, infoGetter info.Getter) []files.Filer {
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

// OutputWriteCloser _
type OutputWriteCloser interface {
	io.WriteCloser
	io.StringWriter
}

// BuildOutputPipe _
func BuildOutputPipe(outputFilePath string) (OutputWriteCloser, error) {
	if outputFilePath != "" {
		outputFile := files.NewFile(outputFilePath)

		if err := outputFile.Create(); err != nil {
			return nil, errors.Wrap(err, "Creating output file")
		}

		outputWriter, err := outputFile.WriteContent()

		if err != nil {
			return nil, errors.Wrap(err, "Initializing output writer")
		}

		return outputWriter, nil
	}

	return os.Stdout, nil
}
