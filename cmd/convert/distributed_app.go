package convert

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxinterrupt"
	"github.com/wailorman/fftb/pkg/distributed/dconfig"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

type distributedApplication struct {
	config      *dconfig.Instance
	accessToken string
	logger      *logrus.Entry
	ctx         context.Context
	dealer      pb.Dealer
}

func initDistributedApplication(ctx context.Context) *distributedApplication {
	var err error
	app := &distributedApplication{}

	if app.config, err = dconfig.New(); err != nil {
		panic(errors.Wrap(err, "Initializing config"))
	}

	if app.logger, err = app.config.BuildLogger(); err != nil {
		panic(errors.Wrap(err, "Initializing logger"))
	}

	if app.accessToken, err = app.config.BuildAccessToken(); err != nil {
		app.logger.WithError(err).Fatal("Building access token")
	}

	app.ctx = ctxinterrupt.ContextWithInterruptHandling(ctx)
	app.dealer = app.config.BuildDealer(app.logger)

	return app
}
