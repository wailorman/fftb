package convert

import (
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/media/ff"
)

// MetadataReceivedBatchMessage _
type MetadataReceivedBatchMessage struct {
	Metadata ffmpegModels.Metadata
	Task     Task
}

// InputVideoCodecDetectedBatchMessage _
type InputVideoCodecDetectedBatchMessage struct {
	Codec string
	Task  Task
}

// BatchProgressMessage _
type BatchProgressMessage struct {
	Progress ff.Progressable
	Task     Task
}

// BatchVideoFilteringMessage _
type BatchVideoFilteringMessage struct {
	Message VideoFileFilteringMessage
	Task    Task
}

// BatchErrorMessage _
type BatchErrorMessage struct {
	Err  error
	Task Task
}
