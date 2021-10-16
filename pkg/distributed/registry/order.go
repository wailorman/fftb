package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/ukvs"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// OrderQueueTimeout _
const OrderQueueTimeout = time.Duration(20 * time.Second)

// Order _
type Order struct {
	ObjectType   string     `json:"object_type"`
	ID           string     `json:"id"`
	Kind         string     `json:"kind"`
	Payload      string     `json:"payload"`
	Publisher    string     `json:"publisher"`
	State        string     `json:"state"`
	SegmentsIDs  []string   `json:"segments_ids"`
	RetriesCount int        `json:"retries_count"`
	RetryAt      *time.Time `json:"retry_at"`
}

// ConvertOrderPayload _
type ConvertOrderPayload struct {
	Params  convert.Params `json:"params"`
	InFile  string         `json:"in_file"`
	OutFile string         `json:"out_file"`
}

// FindOrderByID _
func (r *Instance) FindOrderByID(ctx context.Context, id string) (models.IOrder, error) {
	data, err := r.store.Get(fmt.Sprintf("v1/orders/%s", id))

	if err != nil {
		if errors.Is(err, ukvs.ErrNotFound) {
			return nil, models.ErrNotFound
		}

		return nil, errors.Wrap(err, "Accessing store for order")
	}

	modOrder, err := unmarshalOrderModel(data)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling order model")
	}

	return modOrder, nil
}

// // PickOrderFromQueue _
// func (r *Instance) PickOrderFromQueue(fctx context.Context) (models.IOrder, error) {
// 	if !r.orderQueueLock.TryLockTimeout(OrderQueueTimeout) {
// 		return nil, models.ErrLockTimeoutReached
// 	}

// 	defer r.orderQueueLock.Unlock()

// 	return r.SearchOrder(fctx, func(modOrder models.IOrder) bool {
// 		return modOrder.GetState() == models.OrderStateQueued
// 	})
// }

// SearchOrder _
func (r *Instance) SearchOrder(ctx context.Context, check func(models.IOrder) bool) (models.IOrder, error) {
	orders, err := r.searchOrders(ctx, false, check)

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	return orders[0], nil
}

// SearchAllOrders _
func (r *Instance) SearchAllOrders(ctx context.Context, check func(models.IOrder) bool) ([]models.IOrder, error) {
	orders, err := r.searchOrders(ctx, true, check)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// PersistOrder _
func (r *Instance) PersistOrder(ctx context.Context, modOrder models.IOrder) error {
	orderBefore, _ := r.FindOrderByID(ctx, modOrder.GetID())

	dlog.WithOrder(r.logger, modOrder).
		WithField("after", dlog.JSON(modOrder)).
		WithField("before", dlog.JSON(orderBefore)).
		Trace("Persisting order")

	if modOrder == nil {
		return models.ErrMissingOrder
	}

	if validationErr := modOrder.Validate(); validationErr != nil {
		return validationErr
	}

	data, err := marshalOrderModel(modOrder)

	if err != nil {
		return errors.Wrap(err, "Marshaling db order for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/orders/%s", modOrder.GetID()), data)

	if err != nil {
		return errors.Wrap(err, "Persisting order to store")
	}

	return nil
}

func (r *Instance) searchOrders(fctx context.Context, multiple bool, check func(models.IOrder) bool) ([]models.IOrder, error) {
	ffctx, ffcancel := context.WithCancel(fctx)
	defer ffcancel()

	results, failures := r.store.FindAll(ffctx, "v1/orders/*")
	orders := make([]models.IOrder, 0)

	for {
		select {
		case <-r.ctx.Done():
			return orders, r.ctx.Err()

		case <-fctx.Done():
			return orders, fctx.Err()

		case err := <-failures:
			if err != nil {
				return nil, errors.Wrap(err, "Searching for free order")
			}

		case res, ok := <-results:
			if !ok {
				return orders, nil
			}

			modOrder, err := unmarshalOrderModel(res)

			if err != nil {
				r.logger.WithError(err).
					WithField(dlog.KeyStorePayload, string(res)).
					Warn("Unmarshalling order model from store")

				continue
			}

			if check(modOrder) {
				orders = append(orders, modOrder)

				if !multiple {
					return orders, nil
				}
			}

		case <-time.After(SearchTimeout):
			return nil, models.ErrTimeoutReached
		}
	}
}

func unmarshalOrderModel(data []byte) (models.IOrder, error) {
	dbOrder := &Order{}
	err := dbOrder.unmarshal(data)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling")
	}

	return dbOrder.toModel()
}

func marshalOrderModel(modOrder models.IOrder) ([]byte, error) {
	dbOrder := &Order{}
	err := dbOrder.fromModel(modOrder)

	if err != nil {
		return nil, errors.Wrap(err, "Converting from model")
	}

	return dbOrder.marshal()
}

func (dbOrder *Order) unmarshal(data []byte) error {
	return unmarshalObject(data, ObjectTypeOrder, dbOrder)
}

func (dbOrder *Order) marshal() ([]byte, error) {
	return marshalObject(dbOrder)
}

func (dbOrder *Order) toModel() (models.IOrder, error) {
	var modOrder models.IOrder

	switch dbOrder.Kind {
	case models.ConvertV1Type:
		convOrder := &models.ConvertOrder{
			Identity:     dbOrder.ID,
			Type:         dbOrder.Kind,
			State:        dbOrder.State,
			RetriesCount: dbOrder.RetriesCount,
			RetryAt:      dbOrder.RetryAt,
		}

		if dbOrder.Publisher != "" {
			convOrder.Publisher = &models.Author{Name: dbOrder.Publisher}
		}

		modOrder = convOrder
	default:
		return nil, models.ErrUnknownType
	}

	err := deserializeOrderPayload(dbOrder, modOrder)

	if err != nil {
		return nil, errors.Wrap(err, "Deserializing order payload")
	}

	return modOrder, nil
}

func (dbOrder *Order) fromModel(modOrder models.IOrder) error {
	dbOrder.ObjectType = ObjectTypeOrder
	dbOrder.ID = modOrder.GetID()
	dbOrder.Kind = modOrder.GetType()
	dbOrder.State = modOrder.GetState()
	dbOrder.RetriesCount = modOrder.GetRetriesCount()
	dbOrder.RetryAt = modOrder.GetRetryAt()

	if modOrder.GetPublisher() != nil {
		dbOrder.Publisher = modOrder.GetPublisher().GetName()
	}

	err := serializeOrderPayload(modOrder, dbOrder)

	if err != nil {
		return errors.Wrap(err, "Serializing order payload")
	}

	return nil
}

func serializeOrderPayload(modOrder models.IOrder, dbOrder *Order) error {
	convOrder, ok := modOrder.(*models.ConvertOrder)

	if !ok {
		return models.ErrUnknownType
	}

	payload := &ConvertOrderPayload{
		Params:  convOrder.Params,
		InFile:  convOrder.InFile.FullPath(),
		OutFile: convOrder.OutFile.FullPath(),
	}

	bPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	dbOrder.Payload = string(bPayload)

	return nil
}

func deserializeOrderPayload(dbOrder *Order, modOrder models.IOrder) error {
	if dbOrder.Kind != models.ConvertV1Type {
		return models.ErrUnknownType
	}

	convOrder, ok := modOrder.(*models.ConvertOrder)

	if !ok {
		return models.ErrUnknownType
	}

	convPayload := &ConvertOrderPayload{}

	err := json.Unmarshal([]byte(dbOrder.Payload), convPayload)

	if err != nil {
		return err
	}

	convOrder.Params = convPayload.Params
	convOrder.InFile = files.NewFile(convPayload.InFile)
	convOrder.OutFile = files.NewFile(convPayload.OutFile)

	return nil
}
