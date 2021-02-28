package ff

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/files"
	goffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	goffmpegTranscoder "github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// ErrNotInitialized happened when instance wasn't initialized by Init() func
var ErrNotInitialized = errors.New("not initialized")

// ErrAlreadyStarted happened when Start() func calls more than once
var ErrAlreadyStarted = errors.New("already started")

// ErrAlreadyInitialized happened when Init() func calls more than once
var ErrAlreadyInitialized = errors.New("already initialized")

// ErrProcessTimeout happened when ffmpeg does not send any messages more than ProcessTimeout value
var ErrProcessTimeout = errors.New("ffmpeg process timeout")

// ProcessTimeout is maximum time ffmpeg allowed to not send any messages.
// Once this timeout reached, ErrProcessTimeout will happened
var ProcessTimeout = time.Duration(30 * time.Second)

// Instance _
type Instance struct {
	ctx         context.Context
	initialized bool
	started     bool
	wg          *chwg.ChannelledWaitGroup
	inFile      files.Filer
	outFile     files.Filer
	transcoder  *goffmpegTranscoder.Transcoder
}

// New just initializing & configuring instance before start up
func New(ctx context.Context) *Instance {
	return &Instance{
		ctx: ctx,
		wg:  chwg.New(),
	}
}

// Init receives input & output file objects and initializing transcoder.
// Returns an error if transcoder can't initialize
func (c *Instance) Init(inFile, outFile files.Filer) error {
	if c.initialized {
		return ErrAlreadyInitialized
	}

	c.inFile = inFile
	c.outFile = outFile
	c.transcoder = goffmpegTranscoder.New(c.ctx)

	err := c.transcoder.Initialize(inFile.FullPath(), outFile.FullPath())

	if err != nil {
		return err
	}

	c.initialized = true

	return nil
}

// MediaFile returns goffmpeg's MediaFile object for configuring transcoder input & output
func (c *Instance) MediaFile() *goffmpegModels.Mediafile {
	return c.transcoder.MediaFile()
}

// Start starts ffmpeg process & returns 2 channels.
// progress channel will send progress message ~every 1 sec.
// failures channel will send an error object or nil once operation is done.
func (c *Instance) Start() (
	progress chan Progressable,
	failures chan error,
) {
	progress = make(chan Progressable)
	failures = make(chan error)

	c.wg.Add(1)

	go func() {
		defer close(progress)
		defer close(failures)
		defer c.wg.Done()

		if !c.initialized {
			failures <- ErrNotInitialized
			return
		}

		if c.started {
			failures <- ErrAlreadyStarted
			return
		}

		c.started = true

		done := c.transcoder.Run(true)

		_progress := c.transcoder.Output()

		t := time.NewTimer(ProcessTimeout)
		defer t.Stop()

		for {
			select {
			case <-c.ctx.Done():
				c.transcoder.Stop()
				failures <- c.ctx.Err()
				return

			case progressMessage, ok := <-_progress:
				if ok && progressMessage.FramesProcessed != "" {
					t.Reset(ProcessTimeout)

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
					failures <- err
				}

				return

			case <-t.C:
				c.transcoder.Kill()
				failures <- ErrProcessTimeout
				return
			}
		}
	}()

	return progress, failures
}

// Closed returns channel with finished signal
func (c *Instance) Closed() <-chan struct{} {
	return c.wg.Closed()
}
