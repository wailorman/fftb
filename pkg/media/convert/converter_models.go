package convert

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
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

// BatchTask _
type BatchTask struct {
	Parallelism           int    `json:"parallelism" yaml:"parallelism"`
	StopConversionOnError bool   `json:"stop_conversion_on_error" yaml:"stop_conversion_on_error"`
	Tasks                 []Task `json:"tasks" yaml:"tasks"`
}

// Task _
type Task struct {
	ID      string `json:"id" yaml:"id"`
	InFile  string `json:"in_file" yaml:"in_file"`
	OutFile string `json:"out_file" yaml:"out_file"`
	Params  Params `json:"params" yaml:"params"`
}

// Params _
type Params struct {
	VideoCodec       string `json:"video_codec" yaml:"video_codec"`
	HWAccel          string `json:"hw_accel" yaml:"hw_accel"`
	VideoBitRate     string `json:"video_bit_rate" yaml:"video_bit_rate"`
	VideoQuality     int    `json:"video_quality" yaml:"video_quality"`
	Preset           string `json:"preset" yaml:"preset"`
	Scale            string `json:"scale" yaml:"scale"`
	KeyframeInterval int    `json:"keyframe_interval" yaml:"keyframe_interval"`
	Muxer            string `json:"muxer" yaml:"muxer"`
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

// ErrOutputFileExistsOrIsDirectory _
var ErrOutputFileExistsOrIsDirectory = errors.New("Output file exists or is directory")

// ErrVtbQualityNotSupported _
var ErrVtbQualityNotSupported = errors.New("Video quality option is not supported by Apple VideoToolBox")
