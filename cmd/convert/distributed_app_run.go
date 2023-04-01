package convert

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/distributed/dconfig"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/worker"
)

func distributedRun() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		wg := chwg.New()
		app := initDistributedApplication(c.Context)

		for i := 1; i <= app.config.ThreadsCount(); i++ {
			wg.Add(1)

			go func(threadNum int) {
				var workerLogger *logrus.Entry
				workerLogger = app.logger.WithContext(app.ctx)

				if app.config.ThreadsCount() > 1 {
					workerLogger = app.logger.WithField(dlog.KeyThread, threadNum)
				}

				worker := worker.NewWorker(
					app.config.ThreadConfig(&dconfig.ThreadConfigParams{
						Ctx:    app.ctx,
						Dealer: app.dealer,
						Logger: workerLogger,
						Wg:     wg,

						ThreadNumber:  threadNum,
						Authorization: app.accessToken,
					}),
				)

				worker.Start()
				wg.Done()
			}(i)
		}

		wg.Wait()

		return nil
	}
}
