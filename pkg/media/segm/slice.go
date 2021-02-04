package segm

import (
	"fmt"
	"math/rand"

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
func (s *SliceOperation) Init(req SliceRequest) error {
	if s.initialized {
		return ErrAlreadyInitialized
	}

	var err error

	s.inFile = req.InFile
	s.outPath = req.OutPath
	s.keepTimestamps = req.KeepTimestamps
	s.segmentSec = req.SegmentSec

	err = s.createTmpPath()

	if err != nil {
		return errors.Wrap(err, "Create temp path for segments")
	}

	s.ffworker = ff.New()
	err = s.ffworker.Init(req.InFile, s.tmpPath.BuildFile(segmentPrefix+"%03d"+req.InFile.Extension()))

	if err != nil {
		return errors.Wrap(err, "Initializing ffworker")
	}

	mediaFile := s.ffworker.MediaFile()

	// https://askubuntu.com/a/948449
	// https://trac.ffmpeg.org/wiki/Concatenate
	mediaFile.SetMap("0")
	mediaFile.SetVideoCodec("copy")
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetOutputFormat("segment")
	mediaFile.SetSegmentTime(s.segmentSec)
	mediaFile.SetResetTimestamps(!s.keepTimestamps)

	s.initialized = true

	return nil
}

// Run _
func (s *SliceOperation) Run() (
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

		if s.started {
			failed <- ErrAlreadyInitialized
			return
		}

		_progress, _finished, _failed := s.ffworker.Start()

		for {
			select {
			case <-_finished:
				tmpFiles, err := s.tmpPath.Files()

				if err != nil {
					failed <- errors.Wrap(err, "Getting list of segments files")

					err = s.tmpPath.Destroy()

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

				err := s.tmpPath.Destroy()

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
func (s *SliceOperation) Purge() error {
	if s.tmpPath == nil {
		return ErrNotInitialized
	}

	return s.tmpPath.Destroy()
}

// createTmpPath _
func (s *SliceOperation) createTmpPath() error {
	id := fmt.Sprint(rand.Int())
	s.tmpPath = s.outPath.BuildSubpath("_fftb_chunks_" + id)
	return s.tmpPath.Create()
}
