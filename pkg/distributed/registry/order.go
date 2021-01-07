package registry

import (
	"fmt"

	"github.com/wailorman/fftb/pkg/distributed/ukvs"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// Order _
type Order struct {
	ObjectType string `json:"object_type"`
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Payload    string `json:"payload"`
}

// FindOrderByID _
func (r *Instance) FindOrderByID(id string) (models.IOrder, error) {
	data, err := r.store.Get(fmt.Sprintf("v1/orders/%s", id))

	if err != nil {
		if errors.Is(err, ukvs.ErrNotFound) {
			return nil, models.ErrNotFound
		}

		return nil, errors.Wrap(err, "Accessing store for order")
	}

	dbOrder := &Order{}
	err = unmarshalObject(data, OrderObjectType, dbOrder)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling order")
	}

	if dbOrder.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownOrderType
	}

	modOrder := &models.ConvertOrder{}

	// TODO: dealer tasks <- ??? we have FindSegmentsByOrderID

	modOrder.Identity = dbOrder.ID
	modOrder.Type = dbOrder.Kind

	return modOrder, nil
}

// PersistOrder _
func (r *Instance) PersistOrder(order models.IOrder) error {
	if order == nil {
		return models.ErrMissingOrder
	}

	if order.GetType() != models.ConvertV1Type {
		return models.ErrUnknownOrderType
	}

	payloadStr, err := order.GetPayload()

	if err != nil {
		return errors.Wrap(err, "Getting payload json string")
	}

	dbOrder := &Order{
		ID:         order.GetID(),
		ObjectType: OrderObjectType,
		Kind:       order.GetType(),
		Payload:    payloadStr,
	}

	data, err := marshalObject(dbOrder)

	if err != nil {
		return errors.Wrap(err, "Marshaling db order for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/orders/%s", order.GetID()), data)

	if err != nil {
		return errors.Wrap(err, "Persisting order to store")
	}

	return nil

	// convertOrder := &Order{
	// 	ID:      order.GetID(),
	// 	Kind:    order.GetType(),
	// 	Payload: payloadStr,
	// }

	// getResult := r.gdb.First(&Order{}, order.GetID())

	// var result *gorm.DB

	// if errors.Is(getResult.Error, gorm.ErrRecordNotFound) {
	// 	result = r.gdb.Create(convertOrder)
	// } else {
	// 	result = r.gdb.Save(convertOrder)
	// }

	// if result.Error != nil {
	// 	return errors.Wrap(err, "Failed to persist order")
	// }

	// return nil
}
