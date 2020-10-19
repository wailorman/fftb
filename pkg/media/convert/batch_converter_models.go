package convert

import ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"

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
	Progress Progress
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
