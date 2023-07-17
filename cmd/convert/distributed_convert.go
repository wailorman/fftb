package convert

import (
	"github.com/urfave/cli/v2"
)

// TODO: windows prometheus monitoring https://linuxhint.com/install-monitor-windows-os-prometheus/

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
