package minfo

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/info"
)

var framesListSubcommand = &cli.Command{
	Name: "frames-list",
	Action: func(c *cli.Context) error {
		inputFilePath := c.Args().Get(0)

		if inputFilePath == "" {
			return errors.New("Missing file path argument")
		}

		infoGetter := info.New()
		inputFile := files.NewFile(inputFilePath)

		done, frames, failures := infoGetter.GetFramesList(inputFile)

		for {
			select {
			case frame := <-frames:
				jsonBytes, err := json.Marshal(frame)

				if err != nil {
					return errors.Wrap(err, "Marshaling json")
				}

				fmt.Println(string(jsonBytes))
			case err := <-failures:
				return errors.Wrap(err, "Failed to receive frames")
			case <-done:
				return nil
			}
		}
	},
}
