package minfo

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/info"
)

var framesSummarySubcommand = &cli.Command{
	Name: "frames-summary",
	Action: func(c *cli.Context) error {
		inputFilePath := c.Args().Get(0)

		if inputFilePath == "" {
			return errors.New("Missing file path argument")
		}

		infoGetter := info.New()
		inputFile := files.NewFile(inputFilePath)

		summary, err := infoGetter.GetFramesSummary(inputFile)

		jsonBytes, err := json.Marshal(summary)

		if err != nil {
			return errors.Wrap(err, "Marshaling json")
		}

		fmt.Println(string(jsonBytes))

		return nil
	},
}
