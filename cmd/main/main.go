package main

import (
	"log"
	"os"

	"github.com/wailorman/ffchunker/cmd/convert"
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
		Name:    "chunky",
		Version: "v0.2.0",

		Commands: []*cli.Command{
			etime.CliConfig(),
			split.CliConfig(),
			convert.CliConfig(),
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
