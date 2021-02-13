package convert

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Flags: convertParamsFlags(),
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}
					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					// if err != nil {
					// 	return errors.Wrap(err, "Initializing app")
					// }

					// err = app.StartContracter(c)

					// if err != nil {
					// 	return errors.Wrap(err, "Starting contracter")
					// }

					// <-app.Wait()

					// return nil

					err = app.AddTask(c)

					if err != nil {
						return errors.Wrap(err, "Adding task to queue")
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name: "work",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "worker",
					},
				},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.StartWorker()

					if err != nil {
						return errors.Wrap(err, "Starting worker")
					}

					err = app.StartContracter()

					if err != nil {
						return errors.Wrap(err, "Starting contracter")
					}

					<-app.Wait()

					return nil
				},
			},
		},
	}
}
