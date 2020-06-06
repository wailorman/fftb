package main

import (
	"os"
	"time"

	"github.com/wailorman/ffchunker/cmd/convert"
	"github.com/wailorman/ffchunker/cmd/etime"
	"github.com/wailorman/ffchunker/cmd/log"
	"github.com/wailorman/ffchunker/cmd/split"
	"github.com/wailorman/ffchunker/pkg/ctxlog"

	"github.com/urfave/cli/v2"
)

func main() {
	cliApp()
}

func cliApp() {
	app := &cli.App{
		Name:    "chunky",
		Version: "v0.5.0",

		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Sergey Popov",
				Email: "wailorman@gmail.com",
			},
		},

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "verbosity",
				Aliases: []string{"V"},
				Value:   5,
				Usage: "Verbosity level\n" +
					"                                Possible values:\n" +
					"                                0 - quiet mode, only panics\n" +
					"                                1 - fatal errors\n" +
					"                                2 - regular errors\n" +
					"                                3 - warnings\n" +
					"                                4 - info messages (i.e. progress)\n" +
					"                                5 - debug\n" +
					"                                6 - trace ",
			},
		},

		Before: func(c *cli.Context) error {
			log.SetLoggingLevel(c)
			return nil
		},

		Commands: []*cli.Command{
			etime.CliConfig(),
			split.CliConfig(),
			convert.CliConfig(),
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		ctxlog.Logger.Fatal(err)
	}
}
