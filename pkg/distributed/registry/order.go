package registry

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"gorm.io/gorm"
)

// Order _
type Order struct {
	ID      string
	Kind    string
	Payload string
}

// FindOrderByID _
func (r *SqliteRegistry) FindOrderByID(id string) (models.IOrder, error) {
	order := &Order{}
	result := r.gdb.First(order, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}

	if order.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownOrderType
	}

	convertOrder := &models.ConvertOrder{}

	err := json.Unmarshal([]byte(order.Payload), convertOrder)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling payload")
	}

	// TODO: dealer tasks

	convertOrder.Identity = order.ID
	convertOrder.Type = order.Kind

	return convertOrder, nil
}

// PersistOrder _
func (r *SqliteRegistry) PersistOrder(order models.IOrder) error {
	if order.GetType() != models.ConvertV1Type {
		return models.ErrUnknownOrderType
	}

	payloadStr, err := order.GetPayload()

	if err != nil {
		return errors.Wrap(err, "Getting payload json string")
	}

	convertOrder := &Order{
		ID:      order.GetID(),
		Kind:    order.GetType(),
		Payload: payloadStr,
	}

	getResult := r.gdb.First(&Order{}, order.GetID())

	var result *gorm.DB

	if errors.Is(getResult.Error, gorm.ErrRecordNotFound) {
		result = r.gdb.Create(convertOrder)
	} else {
		result = r.gdb.Save(convertOrder)
	}

	if result.Error != nil {
		return errors.Wrap(err, "Failed to persist order")
	}

	return nil
}
