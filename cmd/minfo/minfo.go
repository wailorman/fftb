package minfo

import (
	"github.com/urfave/cli/v2"
)

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "media-info",
		Aliases: []string{"mi"},
		Usage:   "Retuns media information JSON object",
		Subcommands: []*cli.Command{
			basicSubcommand,
			framesListSubcommand,
			framesSummarySubcommand,
		},
	}
}
