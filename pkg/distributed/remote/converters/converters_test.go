package converters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/remote/converters"
	"github.com/wailorman/fftb/pkg/distributed/test/factories"
)

var f *factories.Builder

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	f = factories.NewBuilder()
}

func teardown() {}

func Test__DealerRequest(t *testing.T) {
	t.Run("DealerRequest: model -> rpc", func(t *testing.T) {
		actual, err := converters.ToRPCDealerRequest(f.Authorization, f.ConvertDealerRequest())

		if assert.NoError(t, err) {
			assert.Equal(t, f.RPCConvertDealerRequest(), actual)
		}
	})

	t.Run("DealerRequest: rpc -> model", func(t *testing.T) {
		actualAuthorization, actualRequest, err := converters.FromRPCDealerRequest(f.RPCConvertDealerRequest())

		if assert.NoError(t, err) {
			assert.Equal(t, f.Authorization, actualAuthorization)
			assert.Equal(t, f.ConvertDealerRequest(), actualRequest)
		}
	})
}

func Test__Segment(t *testing.T) {
	t.Run("Segment: model -> rpc", func(t *testing.T) {
		actual, err := converters.ToRPCSegment(f.ConvertSegment())

		if assert.NoError(t, err) {
			assert.Equal(t, actual, f.RPCConvertSegment())
		}
	})

	t.Run("Segment: rpc -> model", func(t *testing.T) {
		actual, err := converters.FromRPCSegment(f.RPCConvertSegment())

		if assert.NoError(t, err) {
			assert.Equal(t, actual, f.ConvertSegment())
		}
	})
}
