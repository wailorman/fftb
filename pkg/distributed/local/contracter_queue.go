package local

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// AddOrderToQueue _
func (c *ContracterInstance) AddOrderToQueue(ctx context.Context, req models.IContracterRequest) (models.IOrder, error) {
	if validationErr := req.Validate(); validationErr != nil {
		return nil, validationErr
	}

	convertRequest, ok := req.(*models.ConvertContracterRequest)

	if !ok {
		return nil, errors.Wrap(models.ErrUnknownType, fmt.Sprintf("Received request with type `%s`", req.GetType()))
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
		InFile:    convertRequest.InFile,
		OutFile:   convertRequest.OutFile,
	}

	err := c.registry.PersistOrder(ctx, order)

	if err != nil {
		return nil, errors.Wrap(err, "Persisting order")
	}

	return order, nil
}
