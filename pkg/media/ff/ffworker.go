package ff

import (
	"github.com/wailorman/fftb/pkg/files"
	goffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	goffmpegTranscoder "github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// Instance _
type Instance struct {
	Started  chan bool
	Stopping chan bool
	Stopped  chan bool

	stop       chan struct{}
	inFile     files.Filer
	outFile    files.Filer
	transcoder *goffmpegTranscoder.Transcoder
}

// New just initializing & configuring instance before start up
func New() *Instance {
	return &Instance{
		stop: make(chan struct{}),
	}
}

// Init receives input & output file objects and initializing transcoder.
// Returns an error if transcoder can't initialize
func (c *Instance) Init(inFile, outFile files.Filer) error {
	c.inFile = inFile
	c.outFile = outFile
	c.transcoder = new(goffmpegTranscoder.Transcoder)

	err := c.transcoder.Initialize(inFile.FullPath(), outFile.FullPath())

	return err
}

// MediaFile returns goffmpeg's MediaFile object for configuring transcoder input & output
func (c *Instance) MediaFile() *goffmpegModels.Mediafile {
	return c.transcoder.MediaFile()
}

// Stop stops transcoding proccess & sends two messages to two channels
// in respecive order:
// Stopping
// Stopped
func (c *Instance) Stop() {
	c.stop = make(chan struct{})
	// broadcast to all channel receivers
	close(c.stop)
}

func (c *Instance) initChannels() {
	c.Stopping = make(chan bool, 1)
	c.Stopped = make(chan bool, 1)
	c.Started = make(chan bool, 1)
}

func (c *Instance) closeChannels() {
	close(c.Stopping)
	close(c.Stopped)
	close(c.Started)
}

// Start starts ffmpeg process & returns 3 channels.
// progress channel will send progress message ~every 1 sec.
// finished — once.
// failed channel will send an error object
// if something goes wrong & also send a signal to finished channel.
// Also sends a message to Started channel
func (c *Instance) Start() (
	progress chan Progressable,
	finished chan bool,
	failed chan error,
) {
	progress = make(chan Progressable)
	finished = make(chan bool)
	failed = make(chan error)

	c.initChannels()

	go func() {
		defer close(progress)
		defer close(finished)
		defer close(failed)

		defer c.closeChannels()

		done := c.transcoder.Run(true)

		c.Started <- true

		_progress := c.transcoder.Output()

		for {
			select {
			case <-c.stop:
				c.Stopping <- true
				c.transcoder.Stop()
				c.Stopped <- true
				finished <- true
				return

			case progressMessage := <-_progress:
				if progressMessage.FramesProcessed != "" {
					progress <- &Progress{
						framesProcessed: progressMessage.FramesProcessed,
						currentTime:     progressMessage.CurrentTime,
						currentBitrate:  progressMessage.CurrentBitrate,
						progress:        progressMessage.Progress,
						speed:           progressMessage.Speed,
						fps:             progressMessage.FPS,
						file:            c.inFile,
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
