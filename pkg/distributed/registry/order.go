package registry

import (
	"encoding/json"
	"fmt"

	"github.com/wailorman/fftb/pkg/distributed/ukvs"
	"github.com/wailorman/fftb/pkg/media/convert"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// Order _
type Order struct {
	ObjectType string `json:"object_type"`
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Payload    string `json:"payload"`
	Publisher  string `json:"publisher"`
}

// ConvertOrderPayload _
type ConvertOrderPayload struct {
	Params convert.Params `json:"params"`
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

	err = deserializeOrderPayload(dbOrder, modOrder)

	if err != nil {
		return nil, errors.Wrap(err, "Deserializing order payload")
	}

	modOrder.Identity = dbOrder.ID
	modOrder.Type = dbOrder.Kind

	if dbOrder.Publisher != "" {
		modOrder.Publisher = &models.Author{Name: dbOrder.Publisher}
	}

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

	payloadStr, err := serializeOrderPayload(order)

	if err != nil {
		return errors.Wrap(err, "Serializing order payload")
	}

	dbOrder := &Order{
		ID:         order.GetID(),
		ObjectType: OrderObjectType,
		Kind:       order.GetType(),
		Payload:    payloadStr,
	}

	if order.GetPublisher() != nil {
		dbOrder.Publisher = order.GetPublisher().GetName()
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
}

func serializeOrderPayload(modOrder models.IOrder) (string, error) {
	convOrder, ok := modOrder.(*models.ConvertOrder)

	if !ok {
		return "", models.ErrUnknownOrderType
	}

	payload := &ConvertOrderPayload{
		Params: convOrder.Params,
	}

	bPayload, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	return string(bPayload), nil
}

func deserializeOrderPayload(dbOrder *Order, modOrder models.IOrder) error {
	if dbOrder.Kind != models.ConvertV1Type {
		return models.ErrUnknownOrderType
	}

	convOrder, ok := modOrder.(*models.ConvertOrder)

	if !ok {
		return models.ErrUnknownOrderType
	}

	convPayload := &ConvertOrderPayload{}

	err := json.Unmarshal([]byte(dbOrder.Payload), convPayload)

	if err != nil {
		return err
	}

	convOrder.Params = convPayload.Params

	return nil
}
