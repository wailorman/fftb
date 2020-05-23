package main

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/ffchunker"
	"github.com/wailorman/ffchunker/ctxlog"
	"github.com/wailorman/ffchunker/files"

	"github.com/urfave/cli/v2"
)

const bytesInMegabyte = 1000000

func main() {
	cliApp()
}

func cliApp() {
	app := &cli.App{

		Commands: []*cli.Command{
			{
				Name:    "set-times",
				Aliases: []string{"st"},
				Usage:   "Update file modified date meta from it's name",

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "recursively",
						Aliases: []string{"R"},
						Usage:   "Go through all files recursively",
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

					return setTimes(pwd, path, c.Bool("recursively"))
				},
			},
			{
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
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func setTimes(pwd, path string, recursively bool) error {
	if recursively {
		path := files.NewPathBuilder(pwd).NewPath(path)
		resChan, done := ffchunker.NewRecursiveChTimer(path).Perform()

		for {
			select {
			case result := <-resChan:
				logResults(result)
			case <-done:
				return nil
			}
		}
	} else {
		file := files.NewPathBuilder(pwd).NewFile(path)
		res := ffchunker.NewChTimer(file).Perform()

		logResults(res)
	}

	return nil
}

func logResults(result ffchunker.ChTimerResult) {
	log := ctxlog.New(ctxlog.DefaultContext)

	if result.Ok {
		log = log.WithField("time", result.Time.Format(time.RFC3339))
	} else {
		if result.Error != nil {
			log = log.WithField("error", result.Error.Error())
		} else {
			log = log.WithField("error", "No time information")
		}
	}

	log.WithFields(logrus.Fields{
		"full_path":         result.File.FullPath(),
		"used_handler_name": result.UsedHandler,
	}).Info("Time was set")
}

func splitToChunks(pwd, path string, chunkSize int, relativeChunksPath string) error {
	mainFile := files.NewPathBuilder(pwd).NewFile(path)
	outPath := files.NewPathBuilder(mainFile.DirPath()).NewPath(relativeChunksPath)

	log := ctxlog.New(ctxlog.DefaultContext).
		WithFields(logrus.Fields{
			"main_file_path": mainFile.FullPath(),
			"out_path":       outPath.FullPath(),
		})

	log.Info("Splitting to chunks...")

	// (file files.Filer, videoCutter ffchunker.VideoCutter, durationCalculator ffchunker.VideoDurationCalculator, resultPath files.Pather, maxFileSize int) (*ffchunker.Chunker, error)
	chunker, err := ffchunker.NewChunker(
		mainFile,
		ffchunker.NewVideoCutter(),
		ffchunker.NewDurationCalculator(),
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
