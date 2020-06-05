package media

import (
	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/pkg/files"
)

const (
	// NvencHWAccelType _
	NvencHWAccelType = "nvenc"
	// VTBHWAccelType _
	VTBHWAccelType = "videotoolbox"
)

const (
	// HevcCodecType _
	HevcCodecType = "hevc"
	// H264CodecType _
	H264CodecType = "h264"
)

// VideoFileFilteringMessage _
type VideoFileFilteringMessage struct {
	File    files.Filer
	IsVideo bool
	Err     error
}

// ConvertProgress _
type ConvertProgress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
	FPS             int
	File            files.Filer
}

// BatchConverterTask _
type BatchConverterTask struct {
	Parallelism           int
	StopConversionOnError bool
	Tasks                 []ConverterTask
}

// ConverterTask _
type ConverterTask struct {
	ID           string
	InFile       files.Filer
	OutFile      files.Filer
	VideoCodec   string
	HWAccel      string
	VideoBitRate string
	Preset       string
	Scale        string
}

// ErrFileIsNotVideo _
var ErrFileIsNotVideo = errors.New("File is not a video")

// ErrNoStreamsInFile _
var ErrNoStreamsInFile = errors.New("No streams in file")

// ErrCodecIsNotSupportedByEncoder _
var ErrCodecIsNotSupportedByEncoder = errors.New("Codec is not supported by encoder")

// ErrUnsupportedHWAccelType _
var ErrUnsupportedHWAccelType = errors.New("Unsupported hardware acceleration type")

// ErrUnsupportedScale _
var ErrUnsupportedScale = errors.New("Unsupported scale")

// ErrResolutionNotSupportScaling _
var ErrResolutionNotSupportScaling = errors.New("Resolution not support scaling")
