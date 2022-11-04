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

func Test__Segment(t *testing.T) {
	t.Run("model -> rpc", func(t *testing.T) {
		actual, err := converters.ToRPCSegment(f.ConvertSegment())

		if assert.NoError(t, err) {
			assert.Equal(t, actual, f.RPCConvertSegment())
		}
	})

	t.Run("rpc -> model", func(t *testing.T) {
		actual, err := converters.FromRPCSegment(f.RPCConvertSegment())

		if assert.NoError(t, err) {
			assert.Equal(t, actual, f.ConvertSegment())
		}
	})
}
