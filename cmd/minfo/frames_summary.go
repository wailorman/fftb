package minfo

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/minfo"
	"github.com/wailorman/fftb/pkg/media/utils"
)

var framesSummarySubcommand = &cli.Command{
	Name: "frames-summary",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
		},
	},
	Action: func(c *cli.Context) error {
		inputFilePath := c.Args().Get(0)

		if inputFilePath == "" {
			return errors.New("Missing file path argument")
		}

		infoGetter := minfo.New()
		inputFile := files.NewFile(inputFilePath)

		if !inputFile.IsExist() {
			return fmt.Errorf("Input file does not exists: %s", inputFilePath)
		}

		outputWriter, err := utils.BuildOutputPipe(c.String("output"))

		if err != nil {
			return errors.Wrap(err, "Building output pipe")
		}

		summary, err := infoGetter.GetFramesSummary(inputFile)

		jsonBytes, err := json.Marshal(summary)

		if err != nil {
			return errors.Wrap(err, "Marshaling json")
		}

		outputWriter.WriteString(string(jsonBytes))
		outputWriter.WriteString("\n")

		return nil
	},
}
