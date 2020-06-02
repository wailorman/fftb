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
				Name:    "preset",
				Aliases: []string{"p"},
				Usage:   "Encoding preset",
				Value:   "slow",
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

			var progressChan chan media.ConvertProgress
			var doneChan chan bool
			var errChan chan error

			converter := media.NewConverter(mediaInfoGetter)

			if c.Bool("recursively") {
				inputPath := files.NewPath(inputPath)
				outputPath := inputPath

				progressChan, doneChan, errChan = converter.RecursiveConvert(
					media.RecursiveConverterTask{
						InPath:       inputPath,
						OutPath:      outputPath,
						HWAccel:      c.String("hwaccel"),
						VideoCodec:   c.String("video_codec"),
						Preset:       c.String("preset"),
						VideoBitRate: c.String("video_bitrate"),
					},
				)
			} else {
				inputFile := files.NewFile(inputPath)
				outputFile := inputFile.NewWithSuffix("_out")

				progressChan, doneChan, errChan = converter.Convert(
					media.ConverterTask{
						InFile:       inputFile,
						OutFile:      outputFile,
						HWAccel:      c.String("hwaccel"),
						VideoCodec:   c.String("video_codec"),
						Preset:       c.String("preset"),
						VideoBitRate: c.String("video_bitrate"),
					},
				)
			}

			for {
				select {
				case progress, ok := <-progressChan:
					if !ok {
						log.Warn("Error receiving progress message")
						return nil
					}

					log.WithFields(logrus.Fields{
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
						return nil
					}

					if err != nil {
						return err
					}

					log.Warn("Empty error message received")
					return nil

				case <-doneChan:
					log.Info("Conversion done")
					return nil

				case <-converter.ConversionStartedChan:
					log.Info("Conversion started")

				case inputVideoCodec := <-converter.InputVideoCodecDetectedChan:
					log.WithField("input_video_codec", inputVideoCodec).
						Debug("Input video codec detected")
				}
			}
		},
	}
}
