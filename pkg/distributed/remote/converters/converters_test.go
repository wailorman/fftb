package converters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/converters"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
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
	t.Run("model -> rpc", func(t *testing.T) {
		actual, err := converters.ToRPCDealerRequest(f.Authorization, f.ConvertDealerRequest())

		if assert.NoError(t, err) {
			assert.Equal(t, f.RPCConvertDealerRequest(), actual)
		}
	})

	t.Run("rpc -> model", func(t *testing.T) {
		actualAuthorization, actualRequest, err := converters.FromRPCDealerRequest(f.RPCConvertDealerRequest())

		if assert.NoError(t, err) {
			assert.Equal(t, f.Authorization, actualAuthorization)
			assert.Equal(t, f.ConvertDealerRequest(), actualRequest)
		}
	})
}

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

func Test__Progress(t *testing.T) {
	progress := 0.45

	rpcUploadingInput := f.RPCProgress(pb.ProgressNotification_UPLOADING_INPUT, progress)
	modelUploadingInput := f.Progress(models.UploadingInputStep, progress)

	rpcDownloadingInput := f.RPCProgress(pb.ProgressNotification_DOWNLOADING_INPUT, progress)
	modelDownloadingInput := f.Progress(models.DownloadingInputStep, progress)

	rpcProcessing := f.RPCProgress(pb.ProgressNotification_PROCESSING, progress)
	modelProcessing := f.Progress(models.ProcessingStep, progress)

	rpcUploadingOutput := f.RPCProgress(pb.ProgressNotification_UPLOADING_OUTPUT, progress)
	modelUploadingOutput := f.Progress(models.UploadingOutputStep, progress)

	rpcDownloadingOutput := f.RPCProgress(pb.ProgressNotification_DOWNLOADING_OUTPUT, progress)
	modelDownloadingOutput := f.Progress(models.DownloadingOutputStep, progress)

	t.Run("UPLOADING_INPUT: model -> rpc", func(t *testing.T) {
		expected := rpcUploadingInput
		actual, err := converters.ToRPCProgress(f.Authorization, f.SegmentID, modelUploadingInput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("UPLOADING_INPUT: rpc -> model", func(t *testing.T) {
		expected := modelUploadingInput
		authorization, actual, err := converters.FromRPCProgress(rpcUploadingInput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual, "progress")
			assert.Equal(t, f.Authorization, authorization, "authorization")
		}
	})

	t.Run("DOWNLOADING_INPUT: model -> rpc", func(t *testing.T) {
		expected := rpcDownloadingInput
		actual, err := converters.ToRPCProgress(f.Authorization, f.SegmentID, modelDownloadingInput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("DOWNLOADING_INPUT: rpc -> model", func(t *testing.T) {
		expected := modelDownloadingInput
		authorization, actual, err := converters.FromRPCProgress(rpcDownloadingInput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual, "progress")
			assert.Equal(t, f.Authorization, authorization, "authorization")
		}
	})

	t.Run("PROCESSING: model -> rpc", func(t *testing.T) {
		expected := rpcProcessing
		actual, err := converters.ToRPCProgress(f.Authorization, f.SegmentID, modelProcessing)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("PROCESSING: rpc -> model", func(t *testing.T) {
		expected := modelProcessing
		authorization, actual, err := converters.FromRPCProgress(rpcProcessing)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual, "progress")
			assert.Equal(t, f.Authorization, authorization, "authorization")
		}
	})

	t.Run("UPLOADING_OUTPUT: model -> rpc", func(t *testing.T) {
		expected := rpcUploadingOutput
		actual, err := converters.ToRPCProgress(f.Authorization, f.SegmentID, modelUploadingOutput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("UPLOADING_OUTPUT: rpc -> model", func(t *testing.T) {
		expected := modelUploadingOutput
		authorization, actual, err := converters.FromRPCProgress(rpcUploadingOutput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual, "progress")
			assert.Equal(t, f.Authorization, authorization, "authorization")
		}
	})

	t.Run("DOWNLOADING_OUTPUT: model -> rpc", func(t *testing.T) {
		expected := rpcDownloadingOutput
		actual, err := converters.ToRPCProgress(f.Authorization, f.SegmentID, modelDownloadingOutput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("DOWNLOADING_OUTPUT: rpc -> model", func(t *testing.T) {
		expected := modelDownloadingInput
		authorization, actual, err := converters.FromRPCProgress(rpcDownloadingInput)

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual, "progress")
			assert.Equal(t, f.Authorization, authorization, "authorization")
		}
	})
}
