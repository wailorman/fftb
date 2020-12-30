package convert

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/info"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
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

	infoGetter     info.Getter
	stopConversion chan struct{}
}

// RecursiveTask _
type RecursiveTask struct {
	Parallelism int
	InPath      files.Pather
	OutPath     files.Pather
	Params      Params
}

// BuildBatchTaskFromRecursive _
func BuildBatchTaskFromRecursive(task RecursiveTask, infoGetter mediaInfo.Getter) (BatchTask, error) {
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
			ID:      strconv.Itoa(i),
			InFile:  file.FullPath(),
			OutFile: outFile.FullPath(),
			Params:  task.Params,
		})
	}

	return batchTask, nil
}
