package convert

import (
	"github.com/urfave/cli/v2"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Subcommands: []*cli.Command{
			{
				Name:   "run",
				Action: distributedRun(),
			},
		},
	}
}
