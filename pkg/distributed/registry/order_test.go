package registry

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

var ordersTestTable = []struct{ order models.IOrder }{
	{
		&models.ConvertOrder{
			Identity:  "some_id_1",
			Type:      models.ConvertV1Type,
			State:     models.OrderStateQueued,
			Publisher: &models.Author{Name: "testing"},
			Params: convert.Params{
				HWAccel:          "nvenc",
				KeyframeInterval: 30,
				Preset:           "slow",
				Scale:            "1/2",
				VideoBitRate:     "10M",
				VideoQuality:     30,
				VideoCodec:       "h264",
			},
		},
	},
}

func Test__Order__Marshalling(t *testing.T) {
	for i, testItem := range ordersTestTable {
		originalOrder := testItem.order

		orderBytes, err := marshalOrderModel(originalOrder)

		assert.Nil(t, err, fmt.Sprintf("item %d: marshalOrderModel", i))

		newOrder, err := unmarshalOrderModel(orderBytes)

		assert.Nil(t, err, fmt.Sprintf("item %d: unmarshalOrderModel", i))

		assert.Equal(t, "", cmp.Diff(newOrder, originalOrder), fmt.Sprintf("item %d: diff", i))
	}
}

func Test__Order__Persisting(t *testing.T) {
	tmpPath, err := files.NewTempFile("fftb", "test__order__persisting.json")
	assert.Nil(t, err, "init store file")
	store, err := localfile.NewClient(context.TODO(), tmpPath.FullPath())
	assert.Nil(t, err, "localfile initialization")

	registry, err := NewRegistry(context.TODO(), store)

	assert.Nil(t, err, "registry initialization")

	for i, testItem := range ordersTestTable {
		originalOrder := testItem.order

		err := registry.PersistOrder(originalOrder)

		assert.Nil(t, err, fmt.Sprintf("item %d: registry.PersistOrder", i))

		newOrder, err := registry.FindOrderByID(originalOrder.GetID())

		assert.Nil(t, err, fmt.Sprintf("item %d: registry.FindOrderByID", i))
		assert.Equal(t, "", cmp.Diff(newOrder, originalOrder), fmt.Sprintf("item %d: diff", i))
	}
}
