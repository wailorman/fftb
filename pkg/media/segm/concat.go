package segm

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
)

// ConcatOperation _
type ConcatOperation struct {
	outFile          files.Filer
	segments         []*Segment
	segmentsListFile files.Filer
	tmpPath          files.Pather
	ffworker         *ff.Instance
	initialized      bool
	started          bool
}

// ConcatRequest _
type ConcatRequest struct {
	OutFile  files.Filer
	Segments []*Segment
}

// NewConcatOperation _
func NewConcatOperation() *ConcatOperation {
	return &ConcatOperation{}
}

// Init _
func (co *ConcatOperation) Init(req ConcatRequest) error {
	if co.initialized {
		return ErrAlreadyInitialized
	}

	var err error

	co.outFile = req.OutFile
	co.segments = req.Segments

	co.tmpPath, err = createTmpSubdir(co.outFile.BuildPath())

	if err != nil {
		return errors.Wrap(err, "Create temp path for segments list file")
	}

	co.segmentsListFile = co.tmpPath.BuildFile("segments.txt")

	err = co.segmentsListFile.Create()

	if err != nil {
		return errors.Wrap(err, "Create temp segments list file")
	}

	writer, err := co.segmentsListFile.WriteContent()

	if err != nil {
		return errors.Wrap(err, "Building temp segments list file writer")
	}

	segmentsListContent := createSegmentsList(co.segments)

	_, err = io.WriteString(writer, segmentsListContent)

	if err != nil && err != io.EOF {
		return errors.Wrap(err, "Writing segments list")
	}

	co.ffworker = ff.New(context.TODO())
	err = co.ffworker.Init(co.segmentsListFile, req.OutFile)

	if err != nil {
		return errors.Wrap(err, "Initializing ffworker")
	}

	mediaFile := co.ffworker.MediaFile()
	mediaFile.SetUnsafe(true)
	mediaFile.SetVideoCodec("copy")
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetInputFormat("concat")

	co.initialized = true

	return nil
}

// Run _
func (co *ConcatOperation) Run() (finished chan struct{}, progress chan ff.Progressable, failures chan error) {
	finished = make(chan struct{})
	progress = make(chan ff.Progressable)
	failures = make(chan error)

	go func() {
		defer close(finished)
		defer close(progress)
		defer close(failures)

		if !co.initialized {
			failures <- ErrNotInitialized
			return
		}

		if co.started {
			failures <- ErrAlreadyInitialized
			return
		}

		_progress, _finished, _failed := co.ffworker.Start()

		for {
			select {
			case <-_finished:
				return

			case failure := <-_failed:
				failures <- failure
				return

			case progressMessage := <-_progress:
				progress <- progressMessage
			}
		}
	}()

	return finished, progress, failures
}

// Prune _
func (co *ConcatOperation) Prune() error {
	if co.segmentsListFile != nil && co.segmentsListFile.IsExist() {
		return co.segmentsListFile.Remove()
	}

	return nil
}
