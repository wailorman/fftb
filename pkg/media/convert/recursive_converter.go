package convert

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	mediaInfo "github.com/wailorman/fftb/pkg/media/minfo"
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
)

// RecursiveConverter _
type RecursiveConverter struct {
	ConversionStarted       chan bool
	TaskConversionStarted   chan Task
	MetadataReceived        chan MetadataReceivedBatchMessage
	InputVideoCodecDetected chan InputVideoCodecDetectedBatchMessage
	ConversionStopping      chan Task
	ConversionStopped       chan Task

	infoGetter     mediaInfo.Getter
	stopConversion chan struct{}
}

// RecursiveConverterTask _
type RecursiveConverterTask struct {
	Parallelism  int
	InPath       files.Pather
	OutPath      files.Pather
	VideoCodec   string
	HWAccel      string
	VideoBitRate string
	VideoQuality int
	Preset       string
	Scale        string
}

// BuildBatchTaskFromRecursive _
func BuildBatchTaskFromRecursive(task RecursiveConverterTask, infoGetter mediaInfo.Getter) (BatchTask, error) {
	allFiles, err := task.InPath.Files()

	if err != nil {
		return BatchTask{}, errors.Wrap(err, "Getting files from path")
	}

	videoFiles := mediaUtils.FilterVideos(allFiles, infoGetter)

	batchTask := BatchTask{
		Parallelism: task.Parallelism,
		Tasks:       make([]Task, 0),
	}

	for i, file := range videoFiles {
		outFile := file.Clone()
		outFile.SetDirPath(task.OutPath)

		batchTask.Tasks = append(batchTask.Tasks, Task{
			ID:           strconv.Itoa(i),
			InFile:       file,
			OutFile:      outFile,
			VideoCodec:   task.VideoCodec,
			HWAccel:      task.HWAccel,
			VideoBitRate: task.VideoBitRate,
			VideoQuality: task.VideoQuality,
			Preset:       task.Preset,
			Scale:        task.Scale,
		})
	}

	return batchTask, nil
}
