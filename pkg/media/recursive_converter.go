package media

import (
	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/pkg/files"
)

// RecursiveConverter _
type RecursiveConverter struct {
	ConversionStarted       chan bool
	TaskConversionStarted   chan ConverterTask
	MetadataReceived        chan MetadataReceivedBatchMessage
	InputVideoCodecDetected chan InputVideoCodecDetectedBatchMessage
	ConversionStopping      chan ConverterTask
	ConversionStopped       chan ConverterTask
	VideoFileFiltered       chan BatchVideoFilteringMessage

	infoGetter     InfoGetter
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
	Preset       string
	Scale        string
}

// BuildBatchTaskFromRecursive _
func BuildBatchTaskFromRecursive(task RecursiveConverterTask, infoGetter InfoGetter) (BatchConverterTask, error) {
	allFiles, err := task.InPath.Files()

	if err != nil {
		return BatchConverterTask{}, errors.Wrap(err, "Getting files from path")
	}

	videoFiles := filterVideos(allFiles, infoGetter)

	batchTask := BatchConverterTask{
		Parallelism: task.Parallelism,
		Tasks:       make([]ConverterTask, 0),
	}

	for _, file := range videoFiles {
		outFile := file.Clone()
		outFile.SetDirPath(task.OutPath)

		batchTask.Tasks = append(batchTask.Tasks, ConverterTask{
			InFile:       file,
			OutFile:      outFile,
			VideoCodec:   task.VideoCodec,
			HWAccel:      task.HWAccel,
			VideoBitRate: task.VideoBitRate,
			Preset:       task.Preset,
			Scale:        task.Scale,
		})
	}

	return batchTask, nil
}
