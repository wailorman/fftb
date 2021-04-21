package remote_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/models"
	mock_models "github.com/wailorman/fftb/pkg/distributed/models/mocks"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/media/convert"
)

type echoClientWrap struct {
	e *echo.Echo
}

func newEchoClientWrap(e *echo.Echo) *echoClientWrap {
	return &echoClientWrap{
		e: e,
	}
}

// Do _
func (ew *echoClientWrap) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	ew.e.ServeHTTP(rec, req)
	return rec.Result(), nil
}

func remotifyDealer(localDealer models.IDealer) (*remote.ClientWithResponses, error) {
	e := echo.New()
	h := handlers.NewDealerHandler(context.Background(), localDealer)

	remote.RegisterHandlers(e, h)

	te := newEchoClientWrap(e)

	cl, err := remote.NewClientWithResponses("http://localhost:8080", remote.WithHTTPClient(te))

	if err != nil {
		return nil, err
	}

	return cl, nil
}

func Test__AllocateSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	// TODO: enrich with parameters
	convertSegRequest := &models.ConvertDealerRequest{
		Type:     models.ConvertV1Type,
		Identity: "123",
		Params:   convert.Params{},
	}

	localDealer.
		EXPECT().
		AllocateSegment(gomock.Any(), gomock.Any(), convertSegRequest).
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient)

		seg, err := rd.AllocateSegment(context.Background(), nil, convertSegRequest)

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}

func Test__GetSegmentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localDealer.
		EXPECT().
		GetSegmentByID(gomock.Any(), gomock.Any(), "123").
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient)

		seg, err := rd.GetSegmentByID(context.Background(), nil, "123")

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}
