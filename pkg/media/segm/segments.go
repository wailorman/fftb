package segm

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
)

const segmentPrefix = "fftb_out_"

// ErrNotInitialized happened when instance wasn't initialized by Init() func
var ErrNotInitialized = errors.New("Segmentor have not been initialized")

// Instance _
type Instance struct {
	inFile         files.Filer
	outPath        files.Pather
	tmpPath        files.Pather
	ffworker       *ff.Instance
	keepTimestamps bool
	segmentSec     int
}

// New _
func New() *Instance {
	return &Instance{}
}

// createTmpPath _
func (s *Instance) createTmpPath() error {
	id := fmt.Sprint(rand.Int())
	s.tmpPath = s.outPath.BuildSubpath(".fftb_chunks" + id)
	return s.tmpPath.Create()
}

// Request _
type Request struct {
	InFile         files.Filer
	OutPath        files.Pather
	KeepTimestamps bool
	SegmentSec     int
}

// Init _
func (s *Instance) Init(req Request) error {
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
	mediaFile.SetMap("0")
	mediaFile.SetVideoCodec("copy")
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetOutputFormat("segment")
	mediaFile.SetSegmentTime(s.segmentSec)
	mediaFile.SetResetTimestamps(!s.keepTimestamps)

	return nil
}

// Segment _
type Segment struct {
	// from 0 to inf
	Position int
	File     files.Filer
}

// Start _
func (s *Instance) Start() (
	progress chan ff.Progressable,
	segments chan *Segment,
	finished chan bool,
	failed chan error,
) {
	progress = make(chan ff.Progressable)
	segments = make(chan *Segment)
	finished = make(chan bool)
	failed = make(chan error)

	go func() {
		_progress, _finished, _failed := s.ffworker.Start()

		defer close(progress)
		defer close(segments)
		defer close(finished)
		defer close(failed)

		for {
			select {
			case <-_finished:
				tmpFiles, err := s.tmpPath.Files()

				if err != nil {
					failed <- errors.Wrap(err, "Getting list of segments files")

					err = s.tmpPath.Destroy()

					if err != nil {
						failed <- errors.Wrap(err, "Removing tmp directory")
					}

					return
				}

				segs := collectSegments(tmpFiles)

				for _, seg := range segs {
					segments <- seg
				}

				finished <- true
				return

			case failure := <-_failed:
				failed <- failure

				err := s.tmpPath.Destroy()

				if err != nil {
					failed <- errors.Wrap(err, "Removing tmp directory")
				}

				return

			case progressMessage := <-_progress:
				progress <- progressMessage
			}
		}
	}()

	return progress, segments, finished, failed
}

// Purge removes all segments from tmp directory & also tmp directory itself
func (s *Instance) Purge() error {
	if s.tmpPath == nil {
		return ErrNotInitialized
	}

	return s.tmpPath.Destroy()
}

func collectSegments(files []files.Filer) []*Segment {
	result := make([]*Segment, 0)

	for _, file := range files {
		foundSegment := getSegmentFromFile(file)

		if foundSegment != nil {
			result = append(result, foundSegment)
		}
	}

	return result
}

func getSegmentFromFile(file files.Filer) *Segment {
	fileName := file.Name()

	reFull := regexp.MustCompile(segmentPrefix + `\d+`)
	reNumber := regexp.MustCompile(`\d+`)

	if !reFull.MatchString(fileName) {
		return nil
	}

	foundStrNum := reNumber.FindString(fileName)

	number, err := strconv.Atoi(foundStrNum)

	if err != nil {
		// TODO:
		fmt.Printf("err: %#v\n", err)
		return nil
	}

	return &Segment{
		Position: number,
		File:     file,
	}
}
