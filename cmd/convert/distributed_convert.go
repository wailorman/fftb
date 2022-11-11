package convert

import (
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/ctxinterrupt"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/distributed/s3"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/files"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Subcommands: []*cli.Command{
			{
				Name: "work",
				Action: func(c *cli.Context) error {
					ctx := ctxinterrupt.ContextWithInterruptHandling(c.Context)

					tmpPath := files.NewTempPath("fftb")
					rpcClient := pb.NewDealerProtobufClient("http://localhost:3000", &http.Client{})
					storageClient := s3.NewStorageClient(tmpPath.FullPath())
					remoteDealer := remote.NewDealer(rpcClient, storageClient)
					w, err := worker.NewWorker(ctx, tmpPath, remoteDealer, storageClient)

					if err != nil {
						return err
					}

					w.Start()

					<-w.Closed()

					return nil
				},
			},
		},
	}
}
