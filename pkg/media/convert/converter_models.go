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

// BatchConverterTask _
type BatchConverterTask struct {
	Parallelism           int             `yaml:"parallelism"`
	StopConversionOnError bool            `yaml:"stop_conversion_on_error"`
	Tasks                 []ConverterTask `yaml:"tasks"`
}

// ConverterTask _
type ConverterTask struct {
	ID               string      `yaml:"id"`
	InFile           files.Filer `yaml:"in_file"`
	OutFile          files.Filer `yaml:"out_file"`
	VideoCodec       string      `yaml:"video_codec"`
	HWAccel          string      `yaml:"hw_accel"`
	VideoBitRate     string      `yaml:"video_bit_rate"`
	VideoQuality     int         `yaml:"video_quality"`
	Preset           string      `yaml:"preset"`
	Scale            string      `yaml:"scale"`
	KeyframeInterval int         `yaml:"keyframe_interval"`
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

// UnmarshalYAML _
func (ct *ConverterTask) UnmarshalYAML(unmarshal func(interface{}) error) error {

	task := struct {
		ID               string `yaml:"id"`
		InFile           string `yaml:"in_file"`
		OutFile          string `yaml:"out_file"`
		VideoCodec       string `yaml:"video_codec"`
		HWAccel          string `yaml:"hw_accel"`
		VideoBitRate     string `yaml:"video_bit_rate"`
		VideoQuality     int    `yaml:"video_quality"`
		Preset           string `yaml:"preset"`
		Scale            string `yaml:"scale"`
		KeyframeInterval int    `yaml:"keyframe_interval"`
	}{}

	if err := unmarshal(&task); err != nil {
		return err
	}

	ct.ID = task.ID
	ct.VideoCodec = task.VideoCodec
	ct.HWAccel = task.HWAccel
	ct.VideoBitRate = task.VideoBitRate
	ct.VideoQuality = task.VideoQuality
	ct.Preset = task.Preset
	ct.Scale = task.Scale
	ct.KeyframeInterval = task.KeyframeInterval

	ct.InFile = files.NewFile(task.InFile)
	ct.OutFile = files.NewFile(task.OutFile)

	return nil
}
