package etime

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/wailorman/ffchunker/pkg/chtime"
	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/files"
)

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "etime",
		Aliases: []string{"et"},
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
	}
}

func setTimes(pwd, path string, recursively bool) error {
	if recursively {
		path := files.NewPathBuilder(pwd).NewPath(path)
		resChan, done := chtime.NewRecursiveChTimer(path).Perform()

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
		res := chtime.NewChTimer(file).Perform()

		logResults(res)
	}

	return nil
}

func logResults(result chtime.ChTimerResult) {
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
