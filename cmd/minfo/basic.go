package minfo

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/info"
)

var basicSubcommand = &cli.Command{
	Name: "basic",
	Action: func(c *cli.Context) error {
		inputFilePath := c.Args().Get(0)

		if inputFilePath == "" {
			return errors.New("Missing file path argument")
		}

		infoGetter := info.New()
		inputFile := files.NewFile(inputFilePath)

		metadata, err := infoGetter.GetMediaInfo(inputFile)

		jsonBytes, err := json.Marshal(metadata)

		if err != nil {
			return errors.Wrap(err, "Marshaling json")
		}

		fmt.Println(string(jsonBytes))

		return nil
	},
}
