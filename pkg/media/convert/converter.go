package convert

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
	mediaInfo "github.com/wailorman/fftb/pkg/media/minfo"
	mediaUtils "github.com/wailorman/fftb/pkg/media/utils"
)

// Converter _
type Converter struct {
	wg         *chwg.ChannelledWaitGroup
	ctx        context.Context
	infoGetter mediaInfo.Getter
	ffworker   *ff.Instance
}

// NewConverter _
func NewConverter(ctx context.Context, infoGetter mediaInfo.Getter) *Converter {
	ffworker := ff.New(ctx)

	return &Converter{
		wg:         chwg.New(),
		ctx:        ctx,
		infoGetter: infoGetter,
		ffworker:   ffworker,
	}
}

// Convert _
func (c *Converter) Convert(task Task) (
	progress chan ff.Progressable,
	failures chan error,
) {
	progress = make(chan ff.Progressable)
	failures = make(chan error)

	c.wg.Add(1)

	go func() {
		var err error

		defer close(progress)
		defer close(failures)
		defer c.wg.Done()

		inFile := files.NewFile(task.InFile)
		outFile := files.NewFile(task.OutFile)

		err = c.ffworker.Init(inFile, outFile)

		if err != nil {
			failures <- errors.Wrap(err, "ffworker initializing error")
			return
		}

		mediaFile := c.ffworker.MediaFile()

		metadata, err := c.infoGetter.GetMediaInfo(inFile)

		if err != nil {
			failures <- errors.Wrap(err, "Getting file metadata")
			return
		}

		// TODO: log metadata

		if !mediaUtils.IsVideo(metadata) {
			failures <- errors.Wrap(err, "Input file is not video")
			return
		}

		codec, err := chooseCodec(task, metadata)

		if err != nil {
			failures <- errors.Wrap(err, "Choosing codec")
			return
		}

		err = codec.configure(mediaFile)

		if err != nil {
			failures <- errors.Wrap(err, "Configuring codec")
			return
		}

		err = newVideoScale(task, metadata).configure(mediaFile)

		if err != nil {
			failures <- errors.Wrap(err, "Configuring video scale")
			return
		}

		err = outFile.BuildPath().Create()

		if err != nil {
			failures <- errors.Wrap(err, "Creating output dir")
			return
		}

		fProgress, fFailures := c.ffworker.Start()

		for {
			select {
			case <-c.ctx.Done():
				failures <- c.ctx.Err()
				return

			case failure, failed := <-fFailures:
				if !failed {
					<-c.ffworker.Closed()
					return
				}

				failures <- failure
				return

			case progressMessage, ok := <-fProgress:
				if ok {
					progress <- progressMessage
				}
			}
		}
	}()

	return progress, failures
}

// Closed _
func (c *Converter) Closed() <-chan struct{} {
	return c.wg.Closed()
}
