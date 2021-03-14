package segm

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
)

const segmentPrefix = "fftb_out_"

// SliceOperation _
type SliceOperation struct {
	ctx            context.Context
	ffctx          context.Context
	ffcancel       func()
	wg             *chwg.ChannelledWaitGroup
	inFile         files.Filer
	outPath        files.Pather
	tmpPath        files.Pather
	ffworker       *ff.Instance
	keepTimestamps bool
	segmentSec     int
	initialized    bool
	started        bool
	logger         logrus.FieldLogger
}

// SliceRequest _
type SliceRequest struct {
	InFile         files.Filer
	OutPath        files.Pather
	KeepTimestamps bool
	SegmentSec     int
}

// NewSliceOperation _
func NewSliceOperation(ctx context.Context) *SliceOperation {
	ffctx, ffcancel := context.WithCancel(ctx)

	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixSegmSliceOperation); logger == nil {
		logger = ctxlog.New(dlog.PrefixSegmSliceOperation)
	}

	return &SliceOperation{
		logger:   logger,
		ffctx:    ffctx,
		ffcancel: ffcancel,
		ctx:      ctx,
		wg:       chwg.New(),
	}
}

// Init _
func (so *SliceOperation) Init(req SliceRequest) error {
	if so.initialized {
		return ErrAlreadyInitialized
	}

	var err error

	so.inFile = req.InFile
	so.outPath = req.OutPath
	so.keepTimestamps = req.KeepTimestamps
	so.segmentSec = req.SegmentSec

	so.tmpPath, err = createTmpSubdir(so.outPath)

	if err != nil {
		return errors.Wrap(err, "Create temp path for segments")
	}

	so.logger = so.logger.WithField("output_path", req.OutPath.FullPath()).
		WithField("input_file", req.InFile.FullPath()).
		WithField("tmp_path", so.tmpPath.FullPath())

	so.ffworker = ff.New(so.ffctx)
	err = so.ffworker.Init(req.InFile, so.tmpPath.BuildFile(segmentPrefix+"%06d"+req.InFile.Extension()))

	if err != nil {
		return errors.Wrap(err, "Initializing ffworker")
	}

	mediaFile := so.ffworker.MediaFile()

	// https://askubuntu.com/a/948449
	// https://trac.ffmpeg.org/wiki/Concatenate
	mediaFile.SetMap("0")
	mediaFile.SetVideoCodec("copy")
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetOutputFormat("segment")
	mediaFile.SetSegmentTime(so.segmentSec)
	mediaFile.SetResetTimestamps(!so.keepTimestamps)

	so.initialized = true

	return nil
}

// Run _
func (so *SliceOperation) Run() (
	progress chan ff.Progressable,
	segments chan *Segment,
	failures chan error,
) {
	progress = make(chan ff.Progressable)
	segments = make(chan *Segment)
	failures = make(chan error)
	so.wg.Add(1)

	so.logger.Debug("Slicing file")

	go func() {
		defer close(progress)
		defer close(segments)
		defer close(failures)
		defer so.wg.Done()

		if !so.initialized {
			failures <- ErrNotInitialized
			return
		}

		if so.started {
			failures <- ErrAlreadyStarted
			return
		}

		fProgress, fFailures := so.ffworker.Start()

		for {
			select {
			case failure, failed := <-fFailures:
				if !failed {
					tmpFiles, err := so.tmpPath.Files()

					if err != nil {
						so.ffcancel()
						<-so.ffworker.Closed()
						failures <- errors.Wrap(err, "Getting list of segments files")
						return
					}

					segs := collectSegments(tmpFiles)

					for _, seg := range segs {
						segments <- seg
					}

					return
				}

				so.ffcancel()
				<-so.ffworker.Closed()
				failures <- failure
				return

			case progressMessage, ok := <-fProgress:
				if ok {
					progress <- progressMessage
				}
			}
		}
	}()

	return progress, segments, failures
}

// Purge _
func (so *SliceOperation) Purge() error {
	var err error

	if so.tmpPath != nil {
		err = so.tmpPath.Destroy()
	}

	if err != nil {
		so.logger.WithError(err).
			Debug("Purging tmp path")
	}

	return err
}

// Closed _
func (so *SliceOperation) Closed() <-chan struct{} {
	return so.wg.Closed()
}
