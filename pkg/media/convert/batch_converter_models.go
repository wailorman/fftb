package convert

import (
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/media/ff"
)

// MetadataReceivedBatchMessage _
type MetadataReceivedBatchMessage struct {
	Metadata ffmpegModels.Metadata
	Task     ConverterTask
}

// InputVideoCodecDetectedBatchMessage _
type InputVideoCodecDetectedBatchMessage struct {
	Codec string
	Task  ConverterTask
}

// BatchProgressMessage _
type BatchProgressMessage struct {
	Progress ff.Progressable
	Task     ConverterTask
}

// BatchVideoFilteringMessage _
type BatchVideoFilteringMessage struct {
	Message VideoFileFilteringMessage
	Task    ConverterTask
}

// BatchErrorMessage _
type BatchErrorMessage struct {
	Err  error
	Task ConverterTask
}
