package main

import (
	"os"
	"time"

	"github.com/wailorman/fftb/cmd/convert"
	"github.com/wailorman/fftb/cmd/etime"
	"github.com/wailorman/fftb/cmd/log"
	"github.com/wailorman/fftb/cmd/split"
	"github.com/wailorman/fftb/pkg/ctxlog"

	"github.com/urfave/cli/v2"
)

func main() {
	cliApp()
}

func cliApp() {
	app := &cli.App{
		Name:    "fftb",
		Version: "v0.7.1",

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
					"\tPossible values:\n" +
					"\t0 - quiet mode, only panics\n" +
					"\t1 - fatal errors\n" +
					"\t2 - regular errors\n" +
					"\t3 - warnings\n" +
					"\t4 - info messages (i.e. progress)\n" +
					"\t5 - debug\n" +
					"\t6 - trace ",
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
