package convert

import (
	"fmt"

	"github.com/wailorman/ffchunker/pkg/files"
	"github.com/wailorman/ffchunker/pkg/media"
	"gopkg.in/yaml.v2"

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
				Value: "slow",
				Usage: "Encoding preset.\n" +
					"\t\n" +
					"\tWARNING! Apple's VideoToolBox is not support presets\t" +
					"\t\n" +
					"\tCPU-encoding values:\n" +
					"\t- ultrafast\n" +
					"\t- superfast\n" +
					"\t- veryfast\n" +
					"\t- faster\n" +
					"\t- fast\n" +
					"\t- medium\n" +
					"\t- slow\n" +
					"\t- slower\n" +
					"\t- veryslow\n" +
					"\t\n" +
					"\tNVENC values:\n" +
					"\t- slow\n" +
					"\t- medium\n" +
					"\t- fast\n" +
					"\t- hp\n" +
					"\t- hq\n" +
					"\t- bd\n" +
					"\t- ll\n" +
					"\t- llhq\n" +
					"\t- llhp\n" +
					"\t- lossless\n" +
					"\t- losslesshp\t",
			},
			&cli.StringFlag{
				Name:  "scale",
				Usage: "Scaling. Possible values: 1/2 (half resolution), 1/4 (quarter resolution)",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Config file path",
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
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Print YAML task file",
			},
		},

		Action: func(c *cli.Context) error {
			mediaInfoGetter := media.NewInfoGetter()

			var progressChan chan media.BatchProgressMessage
			var doneChan chan bool
			var errChan chan media.BatchErrorMessage

			var conversionStarted chan bool
			var inputVideoCodecDetected chan media.InputVideoCodecDetectedBatchMessage

			var batchTask media.BatchConverterTask

			if c.String("config") != "" {
				configFile := files.NewFile(c.String("config"))
				config, err := configFile.ReadContent()

				if err != nil {
					return errors.Wrap(err, "Reading config content")
				}

				err = yaml.Unmarshal([]byte(config), &batchTask)

				if err != nil {
					return errors.Wrap(err, "Parsing config")
				}
			} else {
				inputPath, outputPath, err := pullInputPaths(c)

				if err != nil {
					return errors.Wrap(err, "Getting input & output paths error")
				}

				inFile := files.NewFile(inputPath)
				outFile := inFile.Clone()
				outFile.SetDirPath(files.NewPath(outputPath))

				batchTask = media.BatchConverterTask{
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
			}

			if c.Bool("dry-run") {
				d, err := yaml.Marshal(&batchTask)
				if err != nil {
					return errors.Wrap(err, "Exporting to YAML")
				}

				fmt.Println(string(d))
				return nil
			}

			converter := media.NewBatchConverter(mediaInfoGetter)

			progressChan, doneChan, errChan = converter.Convert(batchTask)

			conversionStarted = converter.ConversionStarted
			inputVideoCodecDetected = converter.InputVideoCodecDetected

			for {
				select {
				case progressMessage := <-progressChan:
					logProgress(progressMessage)

				case errorMessage := <-errChan:
					logError(errorMessage)

				case <-doneChan:
					logDone()
					return nil

				case <-conversionStarted:
					logConversionStarted()

				case msg := <-converter.TaskConversionStarted:
					logTaskConversionStarted(msg)

				case inputVideoCodec := <-inputVideoCodecDetected:
					logInputVideoCodec(inputVideoCodec)
				}
			}
		},
	}
}
