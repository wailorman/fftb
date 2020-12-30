package convert

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/media/ff"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
)

// Converter _
type Converter struct {
	ConversionStarted       chan bool
	MetadataReceived        chan ffmpegModels.Metadata
	InputVideoCodecDetected chan string
	ConversionStopping      chan bool
	ConversionStopped       chan bool

	infoGetter mediaInfo.Getter
	ffworker   *ff.Instance
}

// NewConverter _
func NewConverter(infoGetter mediaInfo.Getter) *Converter {
	ffworker := ff.New()

	return &Converter{
		infoGetter:         infoGetter,
		ffworker:           ffworker,
		ConversionStarted:  ffworker.Started,
		ConversionStopping: ffworker.Stopping,
		ConversionStopped:  ffworker.Stopped,
	}
}

func (c *Converter) initChannels() {
	c.ConversionStopping = make(chan bool, 1)
	c.ConversionStopped = make(chan bool, 1)
	c.ConversionStarted = make(chan bool, 1)
	c.MetadataReceived = make(chan ffmpegModels.Metadata, 1)
	c.InputVideoCodecDetected = make(chan string, 1)
}

func (c *Converter) closeChannels() {
	close(c.ConversionStopping)
	close(c.ConversionStopped)
	close(c.ConversionStarted)
	close(c.MetadataReceived)
	close(c.InputVideoCodecDetected)
}

// Stop _
func (c *Converter) Stop() {
	c.ffworker.Stop()
}

// Convert _
func (c *Converter) Convert(task Task) (
	progress chan ff.Progressable,
	finished chan bool,
	failed chan error,
) {
	progress = make(chan ff.Progressable)
	finished = make(chan bool)
	failed = make(chan error)

	c.initChannels()

	go func() {
		var err error

		defer close(progress)
		defer close(finished)
		defer close(failed)

		defer c.closeChannels()

		inFile := files.NewFile(task.InFile)
		outFile := files.NewFile(task.OutFile)

		err = c.ffworker.Init(inFile, outFile)

		if err != nil {
			failed <- errors.Wrap(err, "ffworker initializing error")
			return
		}

		mediaFile := c.ffworker.MediaFile()

		metadata, err := c.infoGetter.GetMediaInfo(inFile)

		if err != nil {
			failed <- errors.Wrap(err, "Getting file metadata")
			return
		}

		c.MetadataReceived <- metadata

		if !mediaUtils.IsVideo(metadata) {
			failed <- errors.Wrap(err, "Input file is not video")
			return
		}

		c.InputVideoCodecDetected <- mediaUtils.GetVideoCodec(metadata)

		codec, err := chooseCodec(task, metadata)

		if err != nil {
			failed <- errors.Wrap(err, "Choosing codec")
			return
		}

		err = codec.configure(mediaFile)

		if err != nil {
			failed <- errors.Wrap(err, "Configuring codec")
			return
		}

		err = newVideoScale(task, metadata).configure(mediaFile)

		if err != nil {
			failed <- errors.Wrap(err, "Configuring video scale")
			return
		}

		err = outFile.BuildPath().Create()

		if err != nil {
			failed <- errors.Wrap(err, "Creating output dir")
			return
		}

		_progress, _finished, _failed := c.ffworker.Start()

		for {
			select {
			case <-_finished:
				finished <- true
				return

			case failure := <-_failed:
				failed <- failure
				return

			case progressMessage := <-_progress:
				progress <- progressMessage
			}
		}
	}()

	return progress, finished, failed
}
