package convert

import (
	"fmt"

	"github.com/wailorman/fftb/pkg/files"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	mediaConvert "github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/info"
)

// CliConfig _
func CliConfig() *cli.Command {
	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"conv"},
		Usage:   "Convert video",
		UsageText: "single file mode: fftb convert [options] <input file> <output file>\n" +
			"   recursive mode:   fftb convert [options] -R <input path> <output path>\n" +
			"\n" +
			"   If directory does not exists, it will create it for you.\n" +
			"   WARNING: If file already exists, it will overwrite it",
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
				Usage: "Used hardware acceleration type. Possible values:\n" +
					"                                               videotoolbox (for macs),\n" +
					"                                               nvenc (for Nvidia GPUs).\n" +
					"                                               By default uses x264/x265 CPU encoders",
			},
			&cli.StringFlag{
				Name:    "video-bitrate",
				Aliases: []string{"vb"},
				Usage:   "Video bitrate. Ignores if --video-quality is passed. By default delegates choise to ffmpeg. Examples: 25M, 1600K",
			},
			&cli.IntFlag{
				Name:    "video-quality",
				Aliases: []string{"vq"},
				Usage: "Video quality (-crf option for CPU encoding and -qp option for NVENC).\n" +
					"                                      Integer from 1 to 51 (30 is recommended). By default delegates choise to ffmpeg",
			},
			&cli.StringFlag{
				Name:  "scale",
				Usage: "Scaling. Possible values: 1/2 (half resolution), 1/4 (quarter resolution)",
			},
			&cli.IntFlag{
				Name:    "parallelism",
				Aliases: []string{"P"},
				Usage: "Number of parallel ffmpeg workers.\n" +
					"                                  With higher parallelism value you can utilize more CPU/GPU resources, \n" +
					"                                  but in some situations ffmpeg can't run in parallel or will not give a profit",
				Value: 1,
			},
			&cli.BoolFlag{
				Name:    "recursively",
				Aliases: []string{"R"},
				Usage:   "Convert all video files in directory recursively",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Do not execute conversion and print yaml task config",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Config file path (output from --dry-run option)",
			},
			&cli.StringFlag{
				Name:  "preset",
				Value: "slow",
				Usage: "Encoding preset.\n" +
					"\t\n" +
					"\tWARNING! Apple's VideoToolBox does not support presets\n" +
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
		},

		Action: func(c *cli.Context) error {
			infoGetter := info.New()

			var progressChan chan mediaConvert.BatchProgressMessage
			var doneChan chan bool
			var errChan chan mediaConvert.BatchErrorMessage

			var conversionStarted chan bool
			var inputVideoCodecDetected chan mediaConvert.InputVideoCodecDetectedBatchMessage

			var batchTask mediaConvert.BatchTask

			if c.String("config") != "" {
				configFile := files.NewFile(c.String("config"))
				config, err := configFile.ReadAllContent()

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

				batchTask = mediaConvert.BatchTask{
					Parallelism: c.Int("parallelism"),
					Tasks: []mediaConvert.Task{
						mediaConvert.Task{
							InFile:  inFile.FullPath(),
							OutFile: outFile.FullPath(),
							Params: mediaConvert.Params{
								HWAccel:      c.String("hwa"),
								VideoCodec:   c.String("video-codec"),
								Preset:       c.String("preset"),
								VideoBitRate: c.String("video-bitrate"),
								VideoQuality: c.Int("video-quality"),
								Scale:        c.String("scale"),
							},
						},
					},
				}

				if c.Bool("recursively") {
					outputPath := c.Args().Get(1)

					batchTask, err = mediaConvert.BuildBatchTaskFromRecursive(mediaConvert.RecursiveTask{
						Parallelism: c.Int("parallelism"),
						InPath:      files.NewPath(inputPath),
						OutPath:     files.NewPath(outputPath),
						Params: mediaConvert.Params{
							HWAccel:      c.String("hwa"),
							VideoCodec:   c.String("video-codec"),
							Preset:       c.String("preset"),
							VideoBitRate: c.String("video-bitrate"),
							VideoQuality: c.Int("video-quality"),
							Scale:        c.String("scale"),
						},
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

			converter := mediaConvert.NewBatchConverter(infoGetter)

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
