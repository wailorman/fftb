package split

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/chunk"
	"github.com/wailorman/fftb/pkg/media/segm"
)

const bytesInMegabyte = 1000000

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "split",
		Aliases: []string{"sp"},
		Usage:   "Split video file to chunks",
		UsageText: "fftb split [options] <video file path> <output path>\n" +
			"   WARNING! This tool is not tested well and can produce broken files!",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "chunk-size",
				Aliases: []string{"s"},
				Usage:   "Chunk size in seconds (approximate)",
				Value:   60,
			},
		},

		Action: func(c *cli.Context) error {
			ctx := context.Background()
			pwd, err := os.Getwd()

			if err != nil {
				return errors.Wrap(err, "Getting current working directory")
			}

			inputFilePath := c.Args().Get(0)

			if inputFilePath == "" {
				return errors.New("Missing file path argument")
			}

			outputPath := c.Args().Get(1)

			if outputPath == "" {
				return errors.New("Missing output path argument")
			}

			return splitToChunks(ctx, pwd, inputFilePath, c.Int("chunk-size"), outputPath)
		},
	}
}

func splitToChunks(ctx context.Context, pwd, path string, chunkSize int, relativeChunksPath string) error {
	mainFile := files.NewFile(path)
	outPath := files.NewPath(relativeChunksPath)

	log := ctxlog.Logger.
		WithFields(logrus.Fields{
			"main_file_path": mainFile.FullPath(),
			"out_path":       outPath.FullPath(),
		})

	log.Info("Splitting to chunks...")

	segmenter := segm.NewSliceOperation(ctx)
	chunker := chunk.New(ctx, segmenter)
	chunker.Init(chunk.Request{
		InFile:             mainFile,
		OutPath:            outPath,
		SegmentDurationSec: chunkSize,
	})

	cProgress, cFailures := chunker.Start()

	for {
		select {
		case progressMsg := <-cProgress:
			logProgress(progressMsg)
		case failure, failed := <-cFailures:
			if !failed {
				logDone()
				return nil
			}

			logError(failure)
			return nil
		}
	}
}
