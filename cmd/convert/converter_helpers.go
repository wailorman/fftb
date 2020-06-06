package convert

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func pullInputPaths(c *cli.Context) (
	inputPath string,
	outputPath string,
	err error,
) {
	inputPath = c.Args().Get(0)

	if inputPath == "" {
		err = errors.New("Missing input path first argument")
		return
	}

	outputPath = c.Args().Get(1)

	if outputPath == "" {
		err = errors.New("Missing output path second argument")
		return
	}

	return
}
