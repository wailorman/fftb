package ff

import (
	"context"

	"github.com/wailorman/fftb/pkg/files"
	goffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	goffmpegTranscoder "github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// Instance _
type Instance struct {
	closed     chan struct{}
	ctx        context.Context
	inFile     files.Filer
	outFile    files.Filer
	transcoder *goffmpegTranscoder.Transcoder
}

// New just initializing & configuring instance before start up
func New(ctx context.Context) *Instance {
	return &Instance{
		closed: make(chan struct{}),
		ctx:    ctx,
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

	go func() {
		defer close(progress)
		defer close(finished)
		defer close(failed)

		done := c.transcoder.Run(true)

		_progress := c.transcoder.Output()

		for {
			select {
			case <-c.ctx.Done():
				c.transcoder.Stop()
				close(c.closed)
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
				close(c.closed)
				return
			}
		}
	}()

	return progress, finished, failed
}

// Closed _
func (c *Instance) Closed() <-chan struct{} {
	return c.closed
}
