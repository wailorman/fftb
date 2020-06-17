package media

import (
	"github.com/pkg/errors"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
)

// H264Codec _
type H264Codec struct {
	task     ConverterTask
	metadata ffmpegModels.Metadata
}

// NewH264Codec _
func NewH264Codec(task ConverterTask, metadata ffmpegModels.Metadata) *H264Codec {
	return &H264Codec{
		task:     task,
		metadata: metadata,
	}
}

func (c *H264Codec) configure(mediaFile *ffmpegModels.Mediafile) error {
	var err error

	mediaFile.SetVideoCodec("libx264")
	mediaFile.SetPreset(c.task.Preset)
	mediaFile.SetHideBanner(true)
	mediaFile.SetVsync(true)
	mediaFile.SetVideoBitRate(c.task.VideoBitRate)
	mediaFile.SetAudioCodec("copy")

	hwaccel := chooseHwAccel(c.task, c.metadata)

	if err = hwaccel.configure(mediaFile); err != nil {
		return errors.Wrap(err, "Configuring hwaccel")
	}

	return nil
}

func (c *H264Codec) getType() string {
	return H264CodecType
}
