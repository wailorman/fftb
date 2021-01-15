package convert

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/registry"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/files"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	flags := convertParamsFlags()

	flags = append(flags, &cli.BoolFlag{
		Name: "worker",
	})

	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Flags:   flags,
		Action: func(c *cli.Context) error {
			logger := logrus.New()
			logger.SetLevel(logrus.DebugLevel)
			logger.Formatter = new(prefixed.TextFormatter)
			loggerInstance := logger.WithField("prefix", "fftb.distributed")

			ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxlog.LoggerContextKey, loggerInstance))

			storagePath := files.NewPath(".fftb/storage")

			err := storagePath.Create()

			if err != nil {
				cancel()
				panic(err)
			}

			segmentsPath := files.NewPath(".fftb/segments")

			err = segmentsPath.Create()

			if err != nil {
				cancel()
				panic(err)
			}

			workerPath := files.NewPath(".fftb/worker")

			err = workerPath.Create()

			if err != nil {
				cancel()
				panic(err)
			}

			storage := local.NewStorageControl(storagePath)
			store, err := localfile.NewClient(ctx, ".fftb/store.json")

			if err != nil {
				cancel()
				panic(err)
			}

			registry, err := registry.NewRegistry(ctx, store)

			if err != nil {
				cancel()
				panic(err)
			}

			dealer := local.NewDealer(ctx, storage, registry)
			contracter := local.NewContracter(&local.ContracterParameters{
				TempPath: segmentsPath,
				Dealer:   dealer,
			})

			if !c.Bool("worker") {
				inFile := files.NewFile(c.Args().Get(0))

				_, err = contracter.PrepareOrder(&models.ConvertContracterRequest{
					InFile: inFile,
					Params: convertParamsFromFlags(c),
				})

				if err != nil {
					panic(err)
				}
			} else {
				worker := worker.NewWorker(ctx, workerPath, dealer)

				worker.Start()
			}

			cancel()
			<-store.Closed()
			return nil
		},
	}
}
