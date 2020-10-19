package convert

import (
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
)

type vtbHWAccel struct {
	task     ConverterTask
	metadata ffmpegModels.Metadata
}

func (hw *vtbHWAccel) configure(mediaFile *ffmpegModels.Mediafile) error {
	if !mediaUtils.IsVideo(hw.metadata) {
		return ErrFileIsNotVideo
	}

	mediaFile.SetHardwareAcceleration("videotoolbox")
	mediaFile.SetPreset("")

	switch hw.task.VideoCodec {
	case HevcCodecType:
		mediaFile.SetVideoCodec("hevc_videotoolbox")
	case H264CodecType:
		mediaFile.SetVideoCodec("h264_videotoolbox")
	default:
		return ErrCodecIsNotSupportedByEncoder
	}

	return nil
}

func (hw *vtbHWAccel) getType() string {
	return VTBHWAccelType
}
