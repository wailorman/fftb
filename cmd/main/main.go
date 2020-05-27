package main

import (
	"log"
	"os"

	"github.com/wailorman/ffchunker/cmd/etime"
	"github.com/wailorman/ffchunker/cmd/split"

	"github.com/urfave/cli/v2"
)

const bytesInMegabyte = 1000000

func main() {
	cliApp()
}

func cliApp() {
	app := &cli.App{

		Commands: []*cli.Command{
			etime.CliConfig(),
			split.CliConfig(),
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
