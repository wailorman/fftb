package chunk

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
	"github.com/wailorman/fftb/pkg/media/segm"
)

// Instance _
type Instance struct {
	ctx        context.Context
	mainFile   files.Filer
	resultPath files.Pather
	req        Request

	chunks      []files.Filer
	segmenter   Segmenter
	middlewares []Middleware

	durationCalculator DurationCalculator
	timecodeExtractor  TimecodeExtractor
}

// Segmenter _
type Segmenter interface {
	Init(req segm.SliceRequest) error
	Run() (progress chan ff.Progressable, segments chan *segm.Segment, failures chan error)
	Purge() error
	Closed() <-chan struct{}
}

// DurationCalculator _
type DurationCalculator interface {
	CalculateDuration(file files.Filer) (float64, error)
}

// TimecodeExtractor _
type TimecodeExtractor interface {
	GetTimecode() (time.Time, error)
}

// Result _
type Result struct {
	file files.Filer
}

// New _
func New(ctx context.Context, segmenter Segmenter) *Instance {
	return &Instance{
		ctx:         ctx,
		segmenter:   segmenter,
		middlewares: make([]Middleware, 0),
	}
}

// Request _
type Request struct {
	InFile             files.Filer
	OutPath            files.Pather
	SegmentDurationSec int
}

// Middleware _
type Middleware interface {
	RenameSegments(req Request, sortedSegments []*segm.Segment) error
}

// Use _
func (c *Instance) Use(m Middleware) {
	c.middlewares = append(c.middlewares, m)
}

// Init _
func (c *Instance) Init(req Request) error {
	c.segmenter = segm.NewSliceOperation(c.ctx)
	c.req = req

	err := c.segmenter.Init(segm.SliceRequest{
		InFile:         req.InFile,
		OutPath:        req.OutPath,
		KeepTimestamps: false,
		SegmentSec:     req.SegmentDurationSec,
	})

	if err != nil {
		return errors.Wrap(err, "Segmenter initialization")
	}

	return nil
}

// Start _
func (c *Instance) Start() (progress chan ff.Progressable, failures chan error) {
	progress = make(chan ff.Progressable)
	failures = make(chan error)

	sProgress, sSegments, sFailed := c.segmenter.Run()

	segs := make([]*segm.Segment, 0)

	go func() {
		defer close(progress)
		defer close(failures)

		for {
			select {
			case progressMsg, ok := <-sProgress:
				if ok {
					progress <- progressMsg
				}

			case segment, ok := <-sSegments:
				if ok {
					segs = append(segs, segment)
				}

			case failure, failed := <-sFailed:
				if !failed {
					<-c.segmenter.Closed()

					err := persistSegments(c.req, segs)

					if err != nil {
						failures <- errors.Wrap(err, "Failed to persist segments to result path")
						c.segmenter.Purge()
						return
					}

					segs = sortSegments(segs)

					for _, mwr := range c.middlewares {
						err := mwr.RenameSegments(c.req, segs)

						if err != nil {
							failures <- errors.Wrap(err, "Failed to perform middleware")
							c.segmenter.Purge()
							return
						}
					}

					c.segmenter.Purge()

					return
				}

				failures <- failure
			}
		}
	}()

	return progress, failures
}

func sortSegments(segs []*segm.Segment) []*segm.Segment {
	sortedSegments := make([]*segm.Segment, 0)

	for _, seg := range segs {
		sortedSegments = append(sortedSegments, seg)
	}

	sort.SliceStable(segs, func(i, j int) bool {
		return segs[i].Position < segs[j].Position
	})

	return sortedSegments
}

func persistSegments(req Request, segs []*segm.Segment) error {
	for _, seg := range segs {
		segmentNewName := strings.Join([]string{
			req.InFile.BaseName(),
			"_",
			fmt.Sprint(seg.Position),
			req.InFile.Extension(),
		}, "")

		segmentNewFile := req.OutPath.BuildFile(segmentNewName)

		err := seg.File.Move(segmentNewFile.FullPath())

		if err != nil {
			return errors.Wrap(err, "Renaming tmp segment file")
		}
	}

	return nil
}
