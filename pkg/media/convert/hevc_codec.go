package convert

import (
	"github.com/pkg/errors"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
)

// HevcCodec _
type HevcCodec struct {
	task     ConverterTask
	metadata ffmpegModels.Metadata
}

// NewHevcCodec _
func NewHevcCodec(task ConverterTask, metadata ffmpegModels.Metadata) *HevcCodec {
	return &HevcCodec{
		task:     task,
		metadata: metadata,
	}
}

func (c *HevcCodec) configure(mediaFile *ffmpegModels.Mediafile) error {
	var err error

	mediaFile.SetVideoCodec("libx265")
	mediaFile.SetPreset(c.task.Preset)
	mediaFile.SetHideBanner(true)
	mediaFile.SetVsync(true)
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetMaxMuxingQueueSize(2048)

	if c.task.VideoQuality > 0 {
		mediaFile.SetLibx265Params(&ffmpegModels.Libx265Params{CRF: uint32(c.task.VideoQuality)})
	} else {
		mediaFile.SetVideoBitRate(c.task.VideoBitRate)
	}

	hwaccel := chooseHwAccel(c.task, c.metadata)

	if err = hwaccel.configure(mediaFile); err != nil {
		return errors.Wrap(err, "Configuring hwaccel")
	}

	return nil
}

func (c *HevcCodec) getType() string {
	return HevcCodecType
}
