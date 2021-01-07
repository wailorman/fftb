package convert

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"

	// "github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/registry"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/distributed/worker"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "video-codec",
				Aliases: []string{"vc"},
				Usage:   "Video codec. Possible values: h264, hevc",
				Value:   "h264",
			},
			&cli.IntFlag{
				Name:    "video-quality",
				Aliases: []string{"vq"},
				Usage: "Video quality (-crf option for CPU encoding and -qp option for NVENC).\n" +
					"                                      Integer from 1 to 51 (30 is recommended). By default delegates choise to ffmpeg",
			},
			&cli.BoolFlag{
				Name: "worker",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
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

			registry, err := registry.NewRegistry(store)

			if err != nil {
				cancel()
				panic(err)
			}

			dealer := local.NewDealer(storage, registry)
			contracter := local.NewContracter(&local.ContracterParameters{
				TempPath: segmentsPath,
				Dealer:   dealer,
			})

			if !c.Bool("worker") {
				inFile := files.NewFile(c.Args().Get(0))

				_, err = contracter.PrepareOrder(&models.ConvertContracterRequest{
					InFile: inFile,
					Params: convert.Params{
						VideoCodec:   c.String("video-codec"),
						VideoQuality: c.Int("video-quality"),
					},
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
