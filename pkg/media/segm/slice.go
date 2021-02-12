package segm

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
)

const segmentPrefix = "fftb_out_"

// SliceOperation _
type SliceOperation struct {
	inFile         files.Filer
	outPath        files.Pather
	tmpPath        files.Pather
	ffworker       *ff.Instance
	keepTimestamps bool
	segmentSec     int
	initialized    bool
	started        bool
}

// SliceRequest _
type SliceRequest struct {
	InFile         files.Filer
	OutPath        files.Pather
	KeepTimestamps bool
	SegmentSec     int
}

// NewSliceOperation _
func NewSliceOperation() *SliceOperation {
	return &SliceOperation{}
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

	so.ffworker = ff.New()
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
	finished chan struct{},
	progress chan ff.Progressable,
	segments chan *Segment,
	failed chan error,
) {
	finished = make(chan struct{})
	progress = make(chan ff.Progressable)
	segments = make(chan *Segment)
	failed = make(chan error)

	go func() {
		defer close(finished)
		defer close(progress)
		defer close(segments)
		defer close(failed)

		if so.started {
			failed <- ErrAlreadyInitialized
			return
		}

		_progress, _finished, _failed := so.ffworker.Start()

		for {
			select {
			case <-_finished:
				tmpFiles, err := so.tmpPath.Files()

				if err != nil {
					failed <- errors.Wrap(err, "Getting list of segments files")

					err = so.tmpPath.Destroy()

					if err != nil {
						// TODO: log
					}

					return
				}

				segs := collectSegments(tmpFiles)

				for _, seg := range segs {
					segments <- seg
				}

				return

			case failure := <-_failed:
				failed <- failure

				err := so.tmpPath.Destroy()

				if err != nil {
					// TODO: log
				}

				return

			case progressMessage := <-_progress:
				progress <- progressMessage
			}
		}
	}()

	return finished, progress, segments, failed
}

// Purge removes all segments from tmp directory & also tmp directory itself
func (so *SliceOperation) Purge() error {
	if so.tmpPath == nil {
		return ErrNotInitialized
	}

	return so.tmpPath.Destroy()
}
