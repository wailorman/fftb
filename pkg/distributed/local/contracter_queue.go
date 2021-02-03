package local

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// AddOrderToQueue _
func (c *ContracterInstance) AddOrderToQueue(req models.IContracterRequest) (models.IOrder, error) {
	convertRequest, ok := req.(*models.ConvertContracterRequest)

	if !ok {
		return nil, errors.Wrap(models.ErrUnknownRequestType, fmt.Sprintf("Received request with type `%s`", req.GetType()))
	}

	if c.publisher == nil {
		return nil, models.ErrMissingPublisher
	}

	order := &models.ConvertOrder{
		Identity:  uuid.New().String(),
		Type:      models.ConvertV1Type,
		State:     models.OrderStateQueued,
		Params:    convertRequest.Params,
		Publisher: c.publisher,
	}

	err := c.registry.PersistOrder(order)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting order")
	}

	return order, nil
}
