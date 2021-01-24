package convert

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
	mediaInfo "github.com/wailorman/fftb/pkg/media/info"
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
)

// Converter _
type Converter struct {
	closed     chan struct{}
	ctx        context.Context
	infoGetter mediaInfo.Getter
	ffworker   *ff.Instance
}

// NewConverter _
func NewConverter(ctx context.Context, infoGetter mediaInfo.Getter) *Converter {
	ffworker := ff.New(ctx)

	return &Converter{
		closed:     make(chan struct{}),
		ctx:        ctx,
		infoGetter: infoGetter,
		ffworker:   ffworker,
	}
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

	go func() {
		var err error

		defer close(progress)
		defer close(finished)
		defer close(failed)

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

		if !mediaUtils.IsVideo(metadata) {
			failed <- errors.Wrap(err, "Input file is not video")
			return
		}

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
			case <-c.ctx.Done():
				<-c.ffworker.Closed()
				close(c.closed)
				return

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

// Closed _
func (c *Converter) Closed() <-chan struct{} {
	return c.closed
}
