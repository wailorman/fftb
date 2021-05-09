package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dealerSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
)

// NewDealerAPIRouter _
func NewDealerAPIRouter(
	ctx context.Context,
	logger logrus.FieldLogger,
	dealer models.IDealer,
	authoritySecret []byte,
	sessionSecret []byte) *echo.Echo {

	h := NewDealerHandler(ctx, dealer, authoritySecret, sessionSecret)

	e := echo.New()

	e.Use(dlog.EchoLogger(ctxlog.WithPrefix(logger, dlog.PrefixAPI)))
	e.Use(JWTMiddleware(sessionSecret))
	dealerSchema.RegisterHandlers(e, h)

	return e
}
