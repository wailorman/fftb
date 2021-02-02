package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/ukvs"
)

// ErrUnexpectedObjectType _
var ErrUnexpectedObjectType = errors.New("Unexpected Object Type")

// ObjectTypeSegment _
const ObjectTypeSegment = "segment"

// ObjectTypeOrder _
const ObjectTypeOrder = "order"

// // OrdersStorePath _
// const OrdersStorePath = "v1/orders"

// // SegmentsStorePath _
// const SegmentsStorePath = "v1/segments"

// Instance _
type Instance struct {
	ctx             context.Context
	freeSegmentLock trylock.TryLocker
	orderQueueLock  trylock.TryLocker
	store           ukvs.IStore
	logger          logrus.FieldLogger
}

// TypeCheck _
type TypeCheck struct {
	ObjectType string `json:"object_type"`
}

// NewRegistry _
func NewRegistry(ctx context.Context, store ukvs.IStore) (*Instance, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.registry"); logger == nil {
		logger = ctxlog.New("fftb.worker")
	}

	r := &Instance{
		ctx:             ctx,
		freeSegmentLock: trylock.New(),
		orderQueueLock:  trylock.New(),
		store:           store,
		logger:          logger,
	}

	return r, nil
}

// Closed _
func (r *Instance) Closed() <-chan struct{} {
	return r.store.Closed()
}

// Persist writes data to disk
func (r *Instance) Persist() error {
	r.logger.Debug("Persisting registry")

	return r.store.Persist()
}

func unmarshalObject(data []byte, expectedType string, v interface{}) error {
	typedStruct := &TypeCheck{}

	err := json.Unmarshal(data, typedStruct)

	if err != nil {
		return err
	}

	if typedStruct.ObjectType != expectedType {
		return errors.Wrap(ErrUnexpectedObjectType, fmt.Sprintf("Received type `%s`", typedStruct.ObjectType))
	}

	return json.Unmarshal(data, v)
}

func marshalObject(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
