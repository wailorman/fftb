package media

import (
	ffmpegModels "github.com/wailorman/goffmpeg/models"
)

type nvencHWAccel struct {
	task     ConverterTask
	metadata ffmpegModels.Metadata
}

func (hw *nvencHWAccel) configure(mediaFile *ffmpegModels.Mediafile) error {
	if !isVideo(hw.metadata) {
		return ErrFileIsNotVideo
	}

	if len(hw.metadata.Streams) == 0 {
		return ErrNoStreamsInFile
	}

	mediaFile.SetHardwareAcceleration("cuvid")

	switch hw.metadata.Streams[0].CodecName {
	case "hevc":
		mediaFile.SetInputVideoCodec("hevc_cuvid")
	case "h264":
		mediaFile.SetInputVideoCodec("h264_cuvid")
	}

	switch hw.task.VideoCodec {
	case HevcCodecType:
		mediaFile.SetVideoCodec("hevc_nvenc")
		mediaFile.SetVideoTag("hvc1")
	case H264CodecType:
		mediaFile.SetVideoCodec("h264_nvenc")
	default:
		return ErrCodecIsNotSupportedByEncoder
	}

	return nil
}

func (hw *nvencHWAccel) getType() string {
	return NvencHWAccelType
}
