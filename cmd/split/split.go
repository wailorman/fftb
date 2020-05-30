package split

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/files"
	"github.com/wailorman/ffchunker/pkg/media"
)

const bytesInMegabyte = 1000000

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "split",
		Aliases: []string{"sp"},
		Usage:   "Split video file to chunks",

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "chunk-size",
				Aliases: []string{"s"},
				Usage:   "Chunk size in megabytes (approximate)",
				Value:   1024,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Output path for chunks (WARNING: Will overwrite file if it already exists)",
				Value:   "chunks",
			},
		},

		Action: func(c *cli.Context) error {
			pwd, err := os.Getwd()

			if err != nil {
				return errors.Wrap(err, "Getting current working directory")
			}

			path := c.Args().First()

			if path == "" {
				return errors.New("Missing path argument")
			}

			return splitToChunks(pwd, path, c.Int("chunk-size"), c.String("path"))
		},
	}
}

func splitToChunks(pwd, path string, chunkSize int, relativeChunksPath string) error {
	mainFile := files.NewFile(path)
	outPath := files.NewPath(relativeChunksPath)

	log := ctxlog.New(ctxlog.DefaultContext).
		WithFields(logrus.Fields{
			"main_file_path": mainFile.FullPath(),
			"out_path":       outPath.FullPath(),
		})

	log.Info("Splitting to chunks...")

	mediaInfoGetter := media.NewInfoGetter()

	chunker, err := media.NewChunker(
		mainFile,
		media.NewVideoCutter(),
		media.NewDurationCalculator(mediaInfoGetter),
		outPath,
		chunkSize*bytesInMegabyte,
	)

	if err != nil {
		return errors.Wrap(err, "Building chunker")
	}

	err = chunker.Start()

	if err != nil {
		return errors.Wrap(err, "Splitting to chunks")
	}

	log.Info("Splitting done")

	return nil
}
