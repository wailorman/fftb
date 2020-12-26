package convert

import (
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/files"

	// "github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/registry"
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
		},
		Action: func(c *cli.Context) error {
			storagePath := files.NewPath(".fftb/test/storage")

			err := storagePath.Create()

			if err != nil {
				panic(err)
			}

			segmentsPath := files.NewPath(".fftb/test/segments")

			err = segmentsPath.Create()

			if err != nil {
				panic(err)
			}

			storage := local.NewStorageControl(storagePath)
			registry, err := registry.NewSqliteRegistry(".fftb/sqlite.db", "pkg/distributed/registry/migrations/")

			if err != nil {
				panic(err)
			}

			dealer := local.NewDealer(storage, registry)
			contracter := local.NewContracter(&local.ContracterParameters{
				TempPath: segmentsPath,
				Dealer:   dealer,
			})

			inFile := files.NewFile(c.Args().Get(0))

			// order, err := contracter.PrepareOrder(&models.ConvertContracterRequest{
			_, err = contracter.PrepareOrder(&models.ConvertContracterRequest{
				Params: convert.ConverterTask{
					InFile:       inFile,
					VideoCodec:   c.String("video-codec"),
					VideoQuality: c.Int("video-quality"),
				},
			})

			if err != nil {
				panic(err)
			}

			return nil
		},
	}
}
