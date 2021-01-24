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
			pwd, err := os.Getwd()
			ctx := context.Background()

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

	segmenter := segm.New(ctx)
	chunker := chunk.New(ctx, segmenter)
	chunker.Init(chunk.Request{
		InFile:             mainFile,
		OutPath:            outPath,
		SegmentDurationSec: chunkSize,
	})

	progress, finished, failed := chunker.Start()

	for {
		select {
		case progressMsg := <-progress:
			logProgress(progressMsg)
		case failure := <-failed:
			logError(failure)
			return nil
		case <-finished:
			logDone()
			return nil
		}
	}
}
