package media

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/files"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
	"github.com/wailorman/goffmpeg/transcoder"
)

// Converter _
type Converter struct {
	infoGetter                  InfoGetter
	ConversionStartedChan       chan bool
	MetadataReceivedChan        chan ffmpegModels.Metadata
	InputVideoCodecDetectedChan chan string
}

// NewConverter _
func NewConverter(infoGetter InfoGetter) *Converter {
	return &Converter{
		infoGetter: infoGetter,
	}
}

const (
	// NvencHWAccelType _
	NvencHWAccelType = "nvenc"
	// VTBHWAccelType _
	VTBHWAccelType = "videotoolbox"
)

const (
	// HevcCodecType _
	HevcCodecType = "hevc"
	// H264CodecType _
	H264CodecType = "h264"
)

// ConvertProgress _
type ConvertProgress struct {
	FramesProcessed string
	CurrentTime     string
	CurrentBitrate  string
	Progress        float64
	Speed           string
	FPS             int
	File            files.Filer
}

// ConverterTask _
type ConverterTask struct {
	InFile       files.Filer
	OutFile      files.Filer
	VideoCodec   string
	HWAccel      string
	VideoBitRate string
	Preset       string
	Scale        string
}

// RecursiveConverterTask _
type RecursiveConverterTask struct {
	InPath       files.Pather
	OutPath      files.Pather
	VideoCodec   string
	HWAccel      string
	VideoBitRate string
	Preset       string
	Scale        string
}

// ErrFileIsNotVideo _
var ErrFileIsNotVideo = errors.New("File is not a video")

// ErrNoStreamsInFile _
var ErrNoStreamsInFile = errors.New("No streams in file")

// ErrCodecIsNotSupportedByEncoder _
var ErrCodecIsNotSupportedByEncoder = errors.New("Codec is not supported by encoder")

// ErrUnsupportedHWAccelType _
var ErrUnsupportedHWAccelType = errors.New("Unsupported hardware acceleration type")

// ErrUnsupportedScale _
var ErrUnsupportedScale = errors.New("Unsupported scale")

// ErrResolutionNotSupportScaling _
var ErrResolutionNotSupportScaling = errors.New("Resolution not support scaling")

// RecursiveConvert _
func (c *Converter) RecursiveConvert(task RecursiveConverterTask) (
	progressChan chan ConvertProgress,
	doneChan chan bool,
	errChan chan error,
) {
	progressChan = make(chan ConvertProgress)
	doneChan = make(chan bool)
	errChan = make(chan error)

	log := ctxlog.New(ctxlog.DefaultContext + ".recursive-converter")

	go func() {
		defer close(progressChan)
		defer close(doneChan)
		defer close(errChan)

		var err error

		if err = task.OutPath.Create(); err != nil {
			errChan <- errors.Wrap(err, "Creating output path error")
			return
		}

		allFiles, err := task.InPath.Files()

		if err != nil {
			errChan <- errors.Wrap(err, "Listing files in input path error")
			return
		}

		videoFiles := make([]*files.File, 0)

		for _, file := range allFiles {
			fileLog := log.WithField("file_path", file.FullPath())

			mediaInfo, err := c.infoGetter.GetMediaInfo(file)

			if err != nil {
				fileLog.WithField("error", err).Debug("Getting media info error. Ignoring file")
				continue
			}

			if isVideo(mediaInfo) {
				videoFiles = append(videoFiles, file)

				fileLog.Debug("Found video file")
			} else {
				fileLog.Debug("File is not video")
			}
		}

		log.WithField("found_videos_count", len(videoFiles)).
			Debug("Video files filtering done")

		for fileIndex, file := range videoFiles {
			err := c.proceedRecursiveConvert(file, task, progressChan)

			if err != nil {
				errChan <- errors.Wrap(err, "Processing recursive conversion task error")
				return
			}

			log.WithFields(logrus.Fields{
				"file_index":  fileIndex + 1,
				"total_files": len(videoFiles),
				"file_path":   file.FullPath(),
			}).Info("File converted")
		}
	}()

	return progressChan, doneChan, errChan
}

func (c *Converter) proceedRecursiveConvert(file files.Filer, task RecursiveConverterTask, progressChan chan ConvertProgress) error {
	_progressChan, _doneChan, _errChan := c.Convert(ConverterTask{
		InFile:       file,
		OutFile:      file.NewWithSuffix("_out"),
		VideoCodec:   task.VideoCodec,
		HWAccel:      task.HWAccel,
		VideoBitRate: task.VideoBitRate,
		Preset:       task.Preset,
		Scale:        task.Scale,
	})

	for {
		select {
		case progress := <-_progressChan:
			progressChan <- progress
		case err := <-_errChan:
			return err
		case <-_doneChan:
			return nil
		}
	}
}

// Convert _
func (c *Converter) Convert(task ConverterTask) (
	progressChan chan ConvertProgress,
	doneChan chan bool,
	errChan chan error,
) {
	progressChan = make(chan ConvertProgress)
	doneChan = make(chan bool)
	errChan = make(chan error)

	c.ConversionStartedChan = make(chan bool, 1)
	c.MetadataReceivedChan = make(chan ffmpegModels.Metadata, 1)
	c.InputVideoCodecDetectedChan = make(chan string, 1)

	go func() {
		var err error

		defer close(progressChan)
		defer close(doneChan)
		defer close(errChan)

		defer close(c.ConversionStartedChan)
		defer close(c.MetadataReceivedChan)
		defer close(c.InputVideoCodecDetectedChan)

		trans := new(transcoder.Transcoder)

		err = trans.Initialize(
			task.InFile.FullPath(),
			task.OutFile.FullPath(),
		)

		if err != nil {
			errChan <- errors.Wrap(err, "Transcoder initializing error")
			return
		}

		metadata, err := c.infoGetter.GetMediaInfo(task.InFile)

		if err != nil {
			errChan <- errors.Wrap(err, "Getting file metadata")
			return
		}

		c.MetadataReceivedChan <- metadata

		if !isVideo(metadata) {
			errChan <- errors.Wrap(err, "Input file is not video")
			return
		}

		c.InputVideoCodecDetectedChan <- getVideoCodec(metadata)

		codec, err := chooseCodec(task, metadata)

		if err != nil {
			errChan <- errors.Wrap(err, "Choosing codec")
			return
		}

		err = codec.configure(trans.MediaFile())

		if err != nil {
			errChan <- errors.Wrap(err, "Configuring codec")
			return
		}

		err = newVideoScale(task, metadata).configure(trans.MediaFile())

		if err != nil {
			errChan <- errors.Wrap(err, "Configuring video scale")
			return
		}

		done := trans.Run(true)

		c.ConversionStartedChan <- true

		progress := trans.Output()

		for {
			select {
			case progressMessage := <-progress:
				if progressMessage.FramesProcessed != "" {
					progressChan <- ConvertProgress{
						FramesProcessed: progressMessage.FramesProcessed,
						CurrentTime:     progressMessage.CurrentTime,
						CurrentBitrate:  progressMessage.CurrentBitrate,
						Progress:        progressMessage.Progress,
						Speed:           progressMessage.Speed,
						FPS:             progressMessage.FPS,
						File:            task.InFile,
					}
				}
			case err := <-done:
				if err != nil {
					errChan <- err
					return
				}

				doneChan <- true
				return
			}
		}
	}()

	return progressChan, doneChan, errChan
}

func isVideo(metadata ffmpegModels.Metadata) bool {
	if len(metadata.Streams) == 0 {
		return false
	}

	if metadata.Streams[0].BitRate == "" {
		return false
	}

	return true
}

func getVideoCodec(metadata ffmpegModels.Metadata) string {
	if !isVideo(metadata) {
		return ""
	}

	return metadata.Streams[0].CodecName
}
