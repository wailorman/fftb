package convert

import (
	"fmt"

	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
)

const (
	// FixedHalfScaleType _
	FixedHalfScaleType = "1/2"
	// FixedQuarterScaleType _
	FixedQuarterScaleType = "1/4"
)

type videoScale struct {
	task     Task
	metadata ffmpegModels.Metadata
}

func newVideoScale(task Task, metadata ffmpegModels.Metadata) *videoScale {
	return &videoScale{
		task:     task,
		metadata: metadata,
	}
}

func (pv *videoScale) configure(mediaFile *ffmpegModels.Mediafile) error {
	if pv.task.Params.Scale == "" {
		return nil
	}

	var width, height int
	origWidth, origHeight := getVideoResolution(pv.metadata)

	if origWidth <= 4 || origHeight <= 4 {
		return ErrResolutionNotSupportScaling
	}

	switch pv.task.Params.Scale {
	case FixedHalfScaleType:
		width = origWidth / 2
		height = origHeight / 2
	case FixedQuarterScaleType:
		width = origWidth / 4
		height = origHeight / 4
	case "":
		return nil
	default:
		return ErrUnsupportedScale
	}

	if pv.task.Params.HWAccel == NvencHWAccelType {
		mediaFile.SetVideoFilter(
			fmt.Sprintf(
				"scale_cuda=%d:%d",
				width,
				height,
			),
		)
	} else {
		mediaFile.SetVideoFilter(
			fmt.Sprintf(
				"scale=%d:%d",
				width,
				height,
			),
		)
	}

	return nil
}

func getVideoResolution(metadata ffmpegModels.Metadata) (width, height int) {
	if !mediaUtils.IsVideo(metadata) {
		return 0, 0
	}

	return metadata.Streams[0].Width, metadata.Streams[0].Height
}
