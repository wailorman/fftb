package convert

import (
	"fmt"

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

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
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
						fmt.Printf("err: %#v\n", err)
						app.logger.WithError(err).Warn("Failed to add task")
						return errors.Wrap(err, "Adding task to queue")
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name: "work",
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
					}

					// err = app.StartWorker()

					// if err != nil {
					// 	return errors.Wrap(err, "Starting worker")
					// }

					err = app.StartContracter()

					if err != nil {
						return errors.Wrap(err, "Starting contracter")
					}

					err = app.StartAPI()

					if err != nil {
						return errors.Wrap(err, "Starting API")
					}

					<-app.Wait()

					return nil
				},
			},
			{
				Name: "rwork",
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.StartRemoteWorker()

					if err != nil {
						return errors.Wrap(err, "Starting remote worker")
					}

					// err = app.StartWorker()

					// if err != nil {
					// 	return errors.Wrap(err, "Starting worker")
					// }

					// err = app.StartContracter()

					// if err != nil {
					// 	return errors.Wrap(err, "Starting contracter")
					// }

					// err = app.StartAPI()

					// if err != nil {
					// 	return errors.Wrap(err, "Starting API")
					// }

					<-app.Wait()

					return nil
				},
			},
			{
				Name:    "list-orders",
				Aliases: []string{"lo"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "Possible values: queued, in_progress, finished",
					},
				},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
					}

					str, err := app.ListOrders(c)

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name:    "show-order",
				Aliases: []string{"so"},
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
					}

					str, err := app.ShowOrder(c.Args().Get(0))

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name:    "list-segments",
				Aliases: []string{"ls"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "Possible values: prepared, published, finished",
					},
				},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
					}

					str, err := app.ListSegments(c, c.Args().Get(0))

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name:    "cancel-order",
				Aliases: []string{"co"},
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.InitLocal()

					if err != nil {
						return errors.Wrap(err, "Initializing app (local)")
					}

					err = app.CancelOrder(c.Args().Get(0))

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
		},
	}
}
