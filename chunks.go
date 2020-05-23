package ffchunker

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker/ctxlog"

	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/files"
)

// Chunker _
type Chunker struct {
	mainFile           files.Filer
	totalDuration      float64
	videoCutter        VideoCutter
	durationCalculator VideoDurationCalculator
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
	videoCutter VideoCutter,
	durationCalculator VideoDurationCalculator,
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
	log := ctxlog.New(ctxlog.DefaultContext + ".chunker")

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
		resultFile := files.NewPathBuilder(c.resultPath.FullPath()).NewFile("./abc_" + strconv.Itoa(i) + ".mp4")

		// fmt.Printf("resultFile.FullPath(): %#v\n", resultFile.FullPath())

		chunkLog := log.WithFields(logrus.Fields{
			"chunk_file_path": resultFile.FullPath(),
			"chunk_number":    i,
		})

		resultFile.Remove()

		err = resultFile.EnsureParentDirExists()

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Creating parent directory for chunk #%d", i))
		}

		chunkLog.Info(
			fmt.Sprintf("Processing chunk #%d", i),
		)

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

		// fmt.Printf("chunkDuration: %#v\n", chunkDuration)

		c.chunks = append(c.chunks, chunk)

		if err != nil {
			return errors.Wrap(err, "Calculating chunk duration")
		}

		c.currentDuration += chunkDuration
	}

	return nil
}
