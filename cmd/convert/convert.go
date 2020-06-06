package convert

import (
	"github.com/wailorman/ffchunker/pkg/ctxlog"
	"github.com/wailorman/ffchunker/pkg/files"
	"github.com/wailorman/ffchunker/pkg/media"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"conv"},
		Usage:   "Convert video",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "video-codec",
				Aliases: []string{"vc"},
				Usage:   "Video codec. Possible values: h264, hevc",
				Value:   "h264",
			},
			&cli.StringFlag{
				Name:    "hardware-acceleration",
				Aliases: []string{"hwa"},
				Usage:   "Used hardware acceleration type. Possible values: videotoolbox, nvenc",
			},
			&cli.StringFlag{
				Name:    "video-bitrate",
				Aliases: []string{"vb"},
				Usage:   "Video bitrate. By default delegates choise to ffmpeg",
			},
			&cli.StringFlag{
				Name:  "preset",
				Usage: "Encoding preset",
				Value: "slow",
			},
			&cli.StringFlag{
				Name:  "scale",
				Usage: "Scaling. Possible values: 1/2 (half resolution), 1/4 (quarter resolution)",
			},
			&cli.IntFlag{
				Name:    "parallelism",
				Aliases: []string{"P"},
				Usage:   "Number of parallel ffmpeg workers",
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    "recursively",
				Aliases: []string{"R"},
				Usage:   "Go through all files recursively",
			},
		},

		Action: func(c *cli.Context) error {
			var err error

			log := ctxlog.New(ctxlog.DefaultContext)

			inputPath, outputPath, err := pullInputPaths(c)

			if err != nil {
				return errors.Wrap(err, "Getting input & output paths error")
			}

			mediaInfoGetter := media.NewInfoGetter()

			var progressChan chan media.BatchProgressMessage
			var doneChan chan bool
			var errChan chan media.BatchErrorMessage

			var conversionStarted chan bool
			var inputVideoCodecDetected chan media.InputVideoCodecDetectedBatchMessage

			inFile := files.NewFile(inputPath)
			outFile := inFile.Clone()
			outFile.SetDirPath(files.NewPath(outputPath))

			batchTask := media.BatchConverterTask{
				Parallelism: c.Int("parallelism"),
				Tasks: []media.ConverterTask{
					media.ConverterTask{
						InFile:       inFile,
						OutFile:      outFile,
						HWAccel:      c.String("hwaccel"),
						VideoCodec:   c.String("video-codec"),
						Preset:       c.String("preset"),
						VideoBitRate: c.String("video-bitrate"),
						Scale:        c.String("scale"),
					},
				},
			}

			if c.Bool("recursively") {
				outputPath := c.Args().Get(1)

				batchTask, err = media.BuildBatchTaskFromRecursive(media.RecursiveConverterTask{
					Parallelism:  c.Int("parallelism"),
					InPath:       files.NewPath(inputPath),
					OutPath:      files.NewPath(outputPath),
					HWAccel:      c.String("hwaccel"),
					VideoCodec:   c.String("video-codec"),
					Preset:       c.String("preset"),
					VideoBitRate: c.String("video-bitrate"),
					Scale:        c.String("scale"),
				}, mediaInfoGetter)

				if err != nil {
					return errors.Wrap(err, "Building recursive task")
				}
			}

			converter := media.NewBatchConverter(mediaInfoGetter)

			progressChan, doneChan, errChan = converter.Convert(batchTask)

			conversionStarted = converter.ConversionStarted
			inputVideoCodecDetected = converter.InputVideoCodecDetected

			for {
				select {
				case progressMessage := <-progressChan:
					logProgress(log, progressMessage)

				case errorMessage := <-errChan:
					logError(log, errorMessage)

				case <-doneChan:
					logDone(log)
					return nil

				case <-conversionStarted:
					logConversionStarted(log)

				case msg := <-converter.TaskConversionStarted:
					logTaskConversionStarted(log, msg)

				case inputVideoCodec := <-inputVideoCodecDetected:
					logInputVideoCodec(log, inputVideoCodec)
				}
			}
		},
	}
}
