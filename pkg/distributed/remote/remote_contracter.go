package remote

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/wailorman/fftb/pkg/distributed/models"
	cSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/contracter"
)

// ContracterAPIWrapper _
type ContracterAPIWrapper interface {
	SearchOrders(ctx echo.Context) error
	GetOrderByID(ctx echo.Context, orderID cSchema.OrderIDParam) error
	CancelOrderByID(ctx echo.Context, orderID cSchema.OrderIDParam) error
}

// Contracter _
type Contracter struct {
}

// GetOrderByID _
func (rc *Contracter) GetOrderByID(ctx context.Context, id string) (models.IOrder, error) {
	panic("not implemented") // TODO:
}

// GetAllOrders _
func (rc *Contracter) GetAllOrders(ctx context.Context) ([]models.IOrder, error) {
	panic("not implemented") // TODO:
}

// SearchAllOrders _
func (rc *Contracter) SearchAllOrders(ctx context.Context, search models.IOrderSearchCriteria) ([]models.IOrder, error) {
	panic("not implemented") // TODO:
}

// CancelOrderByID _
func (rc *Contracter) CancelOrderByID(ctx context.Context, orderID string, reason string) error {
	panic("not implemented") // TODO:
}
