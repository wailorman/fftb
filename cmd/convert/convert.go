package convert

import (
	"github.com/sirupsen/logrus"
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
				Name:    "video_codec",
				Aliases: []string{"vc"},
				Usage:   "Video codec. Possible values: h264, hevc",
				Value:   "h264",
			},
			&cli.StringFlag{
				Name:    "hwaccel",
				Aliases: []string{"hwa"},
				Usage:   "Used hardware acceleration type. Possible values: videotoolbox, nvenc",
			},
			&cli.StringFlag{
				Name:    "video_bitrate",
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
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    "recursively",
				Aliases: []string{"R"},
				Usage:   "Go through all files recursively",
			},
		},

		Action: func(c *cli.Context) error {
			log := ctxlog.New(ctxlog.DefaultContext)

			inputPath := c.Args().First()

			if inputPath == "" {
				return errors.New("Missing path argument")
			}

			mediaInfoGetter := media.NewInfoGetter()

			log.Info("Converting started")

			var progressChan chan media.BatchProgressMessage
			var doneChan chan bool
			var errChan chan error

			var conversionStarted chan bool
			var inputVideoCodecDetected chan media.InputVideoCodecDetectedBatchMessage

			if c.Bool("recursively") {
				inputPath := files.NewPath(inputPath)
				outputPath := inputPath

				converter := media.NewRecursiveConverter(mediaInfoGetter)

				recursiveTask := media.RecursiveConverterTask{
					Parallelism:  c.Int("parallelism"),
					InPath:       inputPath,
					OutPath:      outputPath,
					HWAccel:      c.String("hwaccel"),
					VideoCodec:   c.String("video_codec"),
					Preset:       c.String("preset"),
					VideoBitRate: c.String("video_bitrate"),
					Scale:        c.String("scale"),
				}

				progressChan, doneChan, errChan = converter.Convert(recursiveTask)

				conversionStarted = converter.ConversionStarted
				inputVideoCodecDetected = converter.InputVideoCodecDetected
			} else {
				inputFile := files.NewFile(inputPath)
				outputFile := inputFile.NewWithSuffix("_out")

				batchTask := media.BatchConverterTask{
					Parallelism: c.Int("parallelism"),
					Tasks: []media.ConverterTask{
						media.ConverterTask{
							InFile:       inputFile,
							OutFile:      outputFile,
							HWAccel:      c.String("hwaccel"),
							VideoCodec:   c.String("video_codec"),
							Preset:       c.String("preset"),
							VideoBitRate: c.String("video_bitrate"),
							Scale:        c.String("scale"),
						},
					},
				}

				converter := media.NewBatchConverter(mediaInfoGetter)

				progressChan, doneChan, errChan = converter.Convert(batchTask)

				conversionStarted = converter.ConversionStarted
				inputVideoCodecDetected = converter.InputVideoCodecDetected
			}

			for {
				select {
				case progressMessage, ok := <-progressChan:
					if !ok {
						log.Warn("Error receiving progress message")
						return nil
					}

					progress := progressMessage.Progress

					log.WithFields(logrus.Fields{
						"id":               progressMessage.Task.ID,
						"frames_processed": progress.FramesProcessed,
						"current_time":     progress.CurrentTime,
						"current_bitrate":  progress.CurrentBitrate,
						"progress":         progress.Progress,
						"speed":            progress.Speed,
						"fps":              progress.FPS,
						"file_path":        progress.File.FullPath(),
					}).Info("Converting progress")

				case err, ok := <-errChan:
					if !ok {
						log.Warn("Error receiving error message")
						break
					}

					if err != nil {
						log.Warn(err)
					}

				case <-doneChan:
					log.Info("Conversion done")
					return nil

				case <-conversionStarted:
					log.Info("Conversion started")

				case inputVideoCodec := <-inputVideoCodecDetected:
					log.WithField("input_video_codec", inputVideoCodec.Codec).
						Debug("Input video codec detected")
				}
			}
		},
	}
}
