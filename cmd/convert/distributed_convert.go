package convert

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxinterrupt"
	"github.com/wailorman/fftb/pkg/distributed/dconfig"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/worker"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Subcommands: []*cli.Command{
			{
				Name: "run",
				Action: func(c *cli.Context) error {
					var err error

					wg := chwg.New()

					var config *dconfig.Instance
					if config, err = dconfig.New(); err != nil {
						panic(errors.Wrap(err, "Initializing config"))
					}

					var logger *logrus.Entry
					if logger, err = config.BuildLogger(); err != nil {
						panic(errors.Wrap(err, "Initializing logger"))
					}

					ctx := ctxinterrupt.ContextWithInterruptHandling(c.Context)
					dealer := config.BuildDealer()

					for i := 1; i <= config.ThreadsCount(); i++ {
						wg.Add(1)

						go func(threadNum int) {
							var workerLogger *logrus.Entry
							workerLogger = logger.WithContext(ctx)

							if config.ThreadsCount() > 1 {
								workerLogger = logger.WithField(dlog.KeyThread, threadNum)
							}

							worker := worker.NewWorker(
								config.ThreadConfig(&dconfig.ThreadConfigParams{
									Ctx:    ctx,
									Dealer: dealer,
									Logger: workerLogger,
									Wg:     wg,

									ThreadNumber: threadNum,
								}),
							)

							worker.Start()
							wg.Done()
						}(i)
					}

					wg.Wait()

					return nil
				},
			},
		},
	}
}
