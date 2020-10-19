package convert

import ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"

type codecConfigurator interface {
	configure(mediaFile *ffmpegModels.Mediafile) error
	getType() string
}

func chooseCodec(task ConverterTask, metadata ffmpegModels.Metadata) (codecConfigurator, error) {
	switch task.VideoCodec {
	case H264CodecType:
		return &H264Codec{
			task:     task,
			metadata: metadata,
		}, nil
	case HevcCodecType:
		return &HevcCodec{
			task:     task,
			metadata: metadata,
		}, nil
	default:
		return nil, ErrCodecIsNotSupportedByEncoder
	}
}
