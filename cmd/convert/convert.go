package convert

import (
	"context"
	"fmt"

	"github.com/wailorman/fftb/pkg/files"
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	mediaConvert "github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/minfo"
)

// CliConfig _
func CliConfig() *cli.Command {
	flags := convertParamsFlags()

	flags = append(
		flags,
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
	)

	return &cli.Command{
		Name:    "convert",
		Aliases: []string{"conv"},
		Usage:   "Convert video",
		UsageText: "single file mode: fftb convert [options] <input file> <output file>\n" +
			"   recursive mode:   fftb convert [options] -R <input path> <output path>\n" +
			"\n" +
			"   If directory does not exists, it will create it for you.\n" +
			"   WARNING: If file already exists, it will overwrite it",

		Flags: flags,

		Action: func(c *cli.Context) error {
			ctx := context.Background()

			infoGetter := minfo.New()

			var progressChan chan mediaConvert.BatchProgressMessage
			var errChan chan mediaConvert.BatchErrorMessage

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
						{
							InFile:  inFile.FullPath(),
							OutFile: outFile.FullPath(),
							Params:  convertParamsFromFlags(c),
						},
					},
				}

				if c.Bool("recursively") {
					outputPath := c.Args().Get(1)

					batchTask, err = mediaConvert.BuildBatchTaskFromRecursive(mediaConvert.RecursiveTask{
						Parallelism: c.Int("parallelism"),
						InPath:      files.NewPath(inputPath),
						OutPath:     files.NewPath(outputPath),
						Params:      convertParamsFromFlags(c),
					}, infoGetter)

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

			converter := mediaConvert.NewBatchConverter(ctx, infoGetter)

			progressChan, errChan = converter.Convert(batchTask)

			for {
				select {
				case progressMessage, ok := <-progressChan:
					if ok {
						logProgress(progressMessage)
					}

				case failure, failed := <-errChan:
					if !failed {
						logDone()
						return nil
					}

					logError(failure)
				}
			}
		},
	}
}
