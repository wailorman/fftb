package etime

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/wailorman/fftb/pkg/chtime"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"
)

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "etime",
		Aliases: []string{"et"},
		Usage:   "Update file modified date meta from it's name",
		UsageText: "fftb etime [options] <input dir or file>\n" +
			"\n" +
			"   Supported timestamp patterns:\n" +
			"   NMS 22-05-2020 21-52-13.mp4\n" +
			"   20180505_170735.mp4\n" +
			"   Far Cry New Dawn 2020.02.12 - 23.03.10.00.DVR.mp4\n" +
			"   Far Cry New Dawn 2020.02.12 - 23.03.10.00.mp4\n" +
			"   2016_05_20_15_31_51-ses.mp4\n",
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
	}
}

func setTimes(pwd, path string, recursively bool) error {
	if recursively {
		path := files.NewPath(path)
		resChan, done := chtime.NewRecursive(path).Perform()

		for {
			select {
			case result := <-resChan:
				logResults(result)
			case <-done:
				return nil
			}
		}
	} else {
		file := files.NewFile(path)
		res := chtime.New(file).Perform()

		logResults(res)
	}

	return nil
}

func logResults(result chtime.Result) {
	log := ctxlog.Logger

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
