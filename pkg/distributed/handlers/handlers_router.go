package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/schema"
)

// NewDealerAPIRouter _
func NewDealerAPIRouter(
	ctx context.Context,
	dealer models.IDealer,
	authoritySecret []byte,
	sessionSecret []byte) *echo.Echo {

	h := NewDealerHandler(ctx, dealer, authoritySecret, sessionSecret)

	e := echo.New()

	e.Use(JWTMiddleware(sessionSecret))
	schema.RegisterHandlers(e, h)

	return e
}
