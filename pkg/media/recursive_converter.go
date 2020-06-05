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

// NewRecursiveConverter _
func NewRecursiveConverter(infoGetter InfoGetter) *RecursiveConverter {
	return &RecursiveConverter{
		infoGetter:     infoGetter,
		stopConversion: make(chan struct{}),
	}
}

// StopConversion _
func (rc *RecursiveConverter) StopConversion() {
	// broadcast to all channel receivers
	close(rc.stopConversion)
}

// Convert _
func (rc *RecursiveConverter) Convert(task RecursiveConverterTask) (
	progress chan BatchProgressMessage,
	finished chan bool,
	failed chan error,
) {
	progress = make(chan BatchProgressMessage)
	finished = make(chan bool)
	failed = make(chan error)

	allFiles, err := task.InPath.Files()

	if err != nil {
		failed <- errors.Wrap(err, "Getting filed from path")
		return
	}

	videoFiles := filterVideos(allFiles, rc.infoGetter)

	converter := NewBatchConverter(rc.infoGetter)

	batchTask := BatchConverterTask{
		Parallelism: task.Parallelism,
		Tasks:       make([]ConverterTask, 0),
	}

	for _, file := range videoFiles {
		batchTask.Tasks = append(batchTask.Tasks, ConverterTask{
			InFile:       file,
			OutFile:      file.NewWithSuffix("_out"),
			VideoCodec:   task.VideoCodec,
			HWAccel:      task.HWAccel,
			VideoBitRate: task.VideoBitRate,
			Preset:       task.Preset,
			Scale:        task.Scale,
		})
	}

	progress, finished, failed = converter.Convert(batchTask)

	rc.ConversionStarted = converter.ConversionStarted
	rc.TaskConversionStarted = converter.TaskConversionStarted
	rc.MetadataReceived = converter.MetadataReceived
	rc.InputVideoCodecDetected = converter.InputVideoCodecDetected
	rc.ConversionStopping = converter.ConversionStopping
	rc.ConversionStopped = converter.ConversionStopped

	return progress, finished, failed
}
