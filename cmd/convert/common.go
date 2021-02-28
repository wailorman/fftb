package convert

import (
	"github.com/urfave/cli/v2"
	mediaConvert "github.com/wailorman/fftb/pkg/media/convert"
)

func convertParamsFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

func convertParamsFromFlags(c *cli.Context) mediaConvert.Params {
	return mediaConvert.Params{
		HWAccel:      c.String("hwa"),
		VideoCodec:   c.String("video-codec"),
		Preset:       c.String("preset"),
		VideoBitRate: c.String("video-bitrate"),
		VideoQuality: c.Int("video-quality"),
		Scale:        c.String("scale"),
	}
}
