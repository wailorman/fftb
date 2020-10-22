package chunk

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"
	mediaCut "github.com/wailorman/fftb/pkg/media/cut"
	mediaDuration "github.com/wailorman/fftb/pkg/media/duration"
)

// Chunker _
type Chunker struct {
	mainFile           files.Filer
	totalDuration      float64
	videoCutter        mediaCut.Cutter
	durationCalculator mediaDuration.Calculator
	resultPath         files.Pather
	maxFileSize        int

	chunks          []files.Filer
	currentDuration float64
}

// ChunkerResult _
type ChunkerResult struct {
	file files.Filer
}

// NewChunker _
func NewChunker(
	file files.Filer,
	videoCutter mediaCut.Cutter,
	durationCalculator mediaDuration.Calculator,
	resultPath files.Pather,
	maxFileSize int,
) (*Chunker, error) {

	duration, err := durationCalculator.Calculate(file)

	if err != nil {
		return nil, errors.Wrap(err, "Calculating file duration")
	}

	return &Chunker{
		mainFile:           file,
		totalDuration:      duration,
		durationCalculator: durationCalculator,
		videoCutter:        videoCutter,
		resultPath:         resultPath,
		maxFileSize:        maxFileSize,
		chunks:             make([]files.Filer, 0),
	}, nil
}

// Start _
func (c *Chunker) Start() error {
	log := ctxlog.Logger

	totalDuration, err := c.durationCalculator.Calculate(c.mainFile)

	if err != nil {
		return errors.Wrap(err, "Calculating duration of main file")
	}

	c.totalDuration = totalDuration

	log.WithFields(logrus.Fields{
		"total_duration": totalDuration,
		"file_path":      c.mainFile.FullPath(),
	}).Info("Start processing file...")

	for i := 0; c.currentDuration < c.totalDuration; i++ {
		resultFile := files.NewFile(
			fmt.Sprintf("./%s_%d%s", c.mainFile.BaseName(), i, c.mainFile.Extension()),
		)

		resultFile.SetDirPath(c.resultPath)

		err := c.resultPath.Create()

		if err != nil {
			return errors.Wrap(err, "Creating result directory")
		}

		chunkLog := log.WithFields(logrus.Fields{
			"chunk_file_path": resultFile.FullPath(),
			"chunk_number":    i,
		})

		resultFile.Remove()

		err = resultFile.EnsureParentDirExists()

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Creating parent directory for chunk #%d", i))
		}

		chunkLog.Info("Processing chunk")

		chunk, err := c.videoCutter.CutVideo(
			c.mainFile,
			resultFile,
			c.currentDuration,
			c.maxFileSize,
		)

		if err != nil {
			return errors.Wrap(err, "Cutting video")
		}

		chunkDuration, err := c.durationCalculator.Calculate(chunk)

		if chunkDuration == 0 {
			err = chunk.Remove()

			if err != nil {
				return errors.Wrap(err, "Empty file removing")
			}

			return nil
		}

		chunkLog = chunkLog.WithField("duration", chunkDuration)
		chunkLog.Info("Cutted chunk")

		c.chunks = append(c.chunks, chunk)

		if err != nil {
			return errors.Wrap(err, "Calculating chunk duration")
		}

		c.currentDuration += chunkDuration
	}

	return nil
}