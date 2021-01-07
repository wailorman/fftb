package registry

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/distributed/ukvs"
)

// ErrUnexpectedObjectType _
var ErrUnexpectedObjectType = errors.New("Unexpected Object Type")

// SegmentObjectType _
const SegmentObjectType = "segment"

// OrderObjectType _
const OrderObjectType = "order"

// // OrdersStorePath _
// const OrdersStorePath = "v1/orders"

// // SegmentsStorePath _
// const SegmentsStorePath = "v1/segments"

// Instance _
type Instance struct {
	freeSegmentLock trylock.TryLocker
	store           ukvs.IStore
}

// TypeCheck _
type TypeCheck struct {
	ObjectType string `json:"object_type"`
}

// NewRegistry _
func NewRegistry(store ukvs.IStore) (*Instance, error) {
	r := &Instance{
		freeSegmentLock: trylock.New(),
		store:           store,
	}

	return r, nil
}

func unmarshalObject(data []byte, expectedType string, v interface{}) error {
	typedStruct := &TypeCheck{}

	err := json.Unmarshal(data, typedStruct)

	if err != nil {
		return err
	}

	if typedStruct.ObjectType != expectedType {
		return errors.Wrap(ErrUnexpectedObjectType, fmt.Sprint("Received type", typedStruct.ObjectType))
	}

	return json.Unmarshal(data, v)
}

func marshalObject(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
