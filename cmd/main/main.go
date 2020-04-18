package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wailorman/ffchunker"
	"github.com/wailorman/ffchunker/files"

	"github.com/urfave/cli/v2"
)

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
						return err
					}

					setTimes(pwd, c.Args().First(), c.Bool("recursively"))

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func setTimes(pwd, path string, recursively bool) {
	if recursively {
		path := files.NewPathBuilder(pwd).NewPath(path)
		resChan, done := ffchunker.NewRecursiveChTimer(path).Perform()

		for {
			select {
			case result := <-resChan:
				logResults(result)
			case <-done:
				return
			}
		}
	} else {
		file := files.NewPathBuilder(pwd).NewFile(path)
		res := ffchunker.NewChTimer(file).Perform()

		logResults(res)
	}
}

func logResults(result ffchunker.ChTimerResult) {
	var comment string

	if result.Ok {
		comment = result.Time.Format(time.RFC3339)
	} else {
		if result.Error != nil {
			comment = result.Error.Error()
		} else {
			comment = "No time information"
		}
	}

	fmt.Printf(
		"%s\t(%s)\t%s\n",
		result.File.FullPath(),
		result.UsedHandler,
		comment,
	)
}
