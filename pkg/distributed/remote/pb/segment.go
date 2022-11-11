package pb

import "github.com/wailorman/fftb/pkg/media/convert"

func (s *Segment) GetID() string {
	return s.Id
}

func (s *Segment) GetConvertSegmentParams() convert.Params {
	return convert.Params{
		VideoCodec:       s.ConvertParams.VideoCodec,
		HWAccel:          s.ConvertParams.HwAccel,
		VideoBitRate:     s.ConvertParams.VideoBitRate,
		VideoQuality:     int(s.ConvertParams.VideoQuality),
		Preset:           s.ConvertParams.Preset,
		Scale:            s.ConvertParams.Scale,
		KeyframeInterval: int(s.ConvertParams.KeyframeInterval),
		Muxer:            s.ConvertParams.Muxer,
	}
}
