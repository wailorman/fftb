package media

import ffmpegModels "github.com/wailorman/goffmpeg/models"

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
	Progress ConvertProgress
	Task     ConverterTask
}

// BatchVideoFilteringMessage _
type BatchVideoFilteringMessage struct {
	Message VideoFileFilteringMessage
	Task    ConverterTask
}
