package convert

import (
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
				Usage:   "Video codec. Possible values: h264, hevc. Default: h264",
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
		},

		Action: func(c *cli.Context) error {
			// pwd, err := os.Getwd()

			// if err != nil {
			// 	return errors.Wrap(err, "Getting current working directory")
			// }

			// path := c.Args().First()

			// if path == "" {
			// 	return errors.New("Missing path argument")
			// }

			// return splitToChunks(pwd, path, c.Int("chunk-size"), c.String("path"))

			return nil
		},
	}
}
