package handlers

// import (
// 	"github.com/labstack/echo/v4"
// 	"github.com/pkg/errors"
// 	"github.com/wailorman/fftb/pkg/distributed/models"
// 	cSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/contracter"
// )

// // ContracterHandler _
// type ContracterHandler struct {
// 	contracter models.IContracter
// }

// // NewContracterHandler _
// func NewContracterHandler(localContracter models.IContracter) *ContracterHandler {
// 	// TODO: handler config

// 	return &ContracterHandler{
// 		contracter: localContracter,
// 	}
// }

// func buildConvertOrder(order models.IOrder) (*cSchema.ConvertOrder, error) {
// 	if order == nil {
// 		return nil, models.ErrMissingOrder
// 	}

// 	convOrder, ok := order.(*models.ConvertOrder)

// 	if !ok {
// 		return nil, errors.Wrapf(models.ErrUnknownType, "Unknown order type `%s`", order.GetType())
// 	}

// 	return &cSchema.ConvertOrder{
// 		Type:   models.ConvertV1Type,
// 		Id:     convOrder.Identity,
// 		State:  convOrder.State,
// 		Input:  convOrder.InFile.FullPath(),
// 		Output: convOrder.OutFile.FullPath(),
// 		Params: cSchema.ConvertParams{
// 			HwAccel:          convOrder.Params.HWAccel,
// 			KeyframeInterval: convOrder.Params.KeyframeInterval,
// 			Preset:           convOrder.Params.Preset,
// 			Scale:            convOrder.Params.Scale,
// 			VideoBitRate:     convOrder.Params.VideoBitRate,
// 			VideoCodec:       convOrder.Params.VideoCodec,
// 			VideoQuality:     convOrder.Params.VideoQuality,
// 		},
// 	}, nil
// }

// // SearchOrders _
// // (GET /orders)
// func (ch *ContracterHandler) SearchOrders(c echo.Context) error {
// 	panic("not implemented") // TODO:
// }

// // GetOrderByID _
// // (GET /orders/{orderID})
// func (ch *ContracterHandler) GetOrderByID(c echo.Context, orderID cSchema.OrderIDParam) error {
// 	order, err := ch.contracter.GetOrderByID(c.Request().Context(), string(orderID))

// 	if err != nil {
// 		return c.JSON(newAPIError(err))
// 	}

// 	responseOrder, err := buildConvertOrder(order)

// 	if err != nil {
// 		return c.JSON(newAPIError(err))
// 	}

// 	return c.JSON(200, responseOrder)
// }

// // CancelOrderByID _
// // (GET /orders/{orderID}/cancel)
// func (ch *ContracterHandler) CancelOrderByID(c echo.Context, orderID cSchema.OrderIDParam) error {
// 	panic("not implemented") // TODO:
// }
