package media

import (
	"github.com/pkg/errors"
	ffmpegModels "github.com/wailorman/goffmpeg/models"
	"github.com/wailorman/goffmpeg/transcoder"
)

// Converter _
type Converter struct {
	ConversionStarted       chan bool
	MetadataReceived        chan ffmpegModels.Metadata
	InputVideoCodecDetected chan string
	ConversionStopping      chan bool
	ConversionStopped       chan bool
	VideoFileFiltered       chan VideoFileFilteringMessage

	infoGetter     InfoGetter
	stopConversion chan struct{}
}

// NewConverter _
func NewConverter(infoGetter InfoGetter) *Converter {
	return &Converter{
		infoGetter:     infoGetter,
		stopConversion: make(chan struct{}),
	}
}

// StopConversion _
func (c *Converter) StopConversion() {
	c.stopConversion = make(chan struct{})
	// broadcast to all channel receivers
	close(c.stopConversion)
}

// initChannels _
func (c *Converter) initChannels() {
	c.ConversionStopping = make(chan bool, 1)
	c.ConversionStopped = make(chan bool, 1)
	c.ConversionStarted = make(chan bool, 1)
	c.MetadataReceived = make(chan ffmpegModels.Metadata, 1)
	c.InputVideoCodecDetected = make(chan string, 1)
}

// closeChannels _
func (c *Converter) closeChannels() {
	close(c.ConversionStopping)
	close(c.ConversionStopped)
	close(c.ConversionStarted)
	close(c.MetadataReceived)
	close(c.InputVideoCodecDetected)
}

// Convert _
func (c *Converter) Convert(task ConverterTask) (
	progress chan ConvertProgress,
	finished chan bool,
	failed chan error,
) {
	progress = make(chan ConvertProgress)
	finished = make(chan bool)
	failed = make(chan error)

	c.initChannels()

	go func() {
		var err error

		defer close(progress)
		defer close(finished)
		defer close(failed)

		defer c.closeChannels()

		trans := new(transcoder.Transcoder)

		err = trans.Initialize(
			task.InFile.FullPath(),
			task.OutFile.FullPath(),
		)

		if err != nil {
			failed <- errors.Wrap(err, "Transcoder initializing error")
			return
		}

		metadata, err := c.infoGetter.GetMediaInfo(task.InFile)

		if err != nil {
			failed <- errors.Wrap(err, "Getting file metadata")
			return
		}

		c.MetadataReceived <- metadata

		if !isVideo(metadata) {
			failed <- errors.Wrap(err, "Input file is not video")
			return
		}

		c.InputVideoCodecDetected <- getVideoCodec(metadata)

		codec, err := chooseCodec(task, metadata)

		if err != nil {
			failed <- errors.Wrap(err, "Choosing codec")
			return
		}

		err = codec.configure(trans.MediaFile())

		if err != nil {
			failed <- errors.Wrap(err, "Configuring codec")
			return
		}

		err = newVideoScale(task, metadata).configure(trans.MediaFile())

		if err != nil {
			failed <- errors.Wrap(err, "Configuring video scale")
			return
		}

		err = task.OutFile.BuildPath().Create()

		if err != nil {
			failed <- errors.Wrap(err, "Creating output dir")
			return
		}

		done := trans.Run(true)

		c.ConversionStarted <- true

		_progress := trans.Output()

		for {
			select {
			case <-c.stopConversion:
				c.ConversionStopping <- true
				trans.Stop()
				c.ConversionStopped <- true
				finished <- true
				return

			case progressMessage := <-_progress:
				if progressMessage.FramesProcessed != "" {
					progress <- ConvertProgress{
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
					failed <- err
				}

				finished <- true
				return
			}
		}
	}()

	return progress, finished, failed
}
