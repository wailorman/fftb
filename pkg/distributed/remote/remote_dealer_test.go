package remote_test

import (
	"context"
	"fmt"
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
	"github.com/wailorman/fftb/pkg/distributed/schema"
	"github.com/wailorman/fftb/pkg/media/convert"
)

var authoritySecret = []byte("authority_secret_remote_dealer")
var sessionSecret = []byte("session_secret_remote_dealer")

type authorEq struct {
	author models.IAuthor
}

// Matches _
func (aeq authorEq) Matches(other interface{}) bool {
	fmt.Printf("other: %#v\n", other)
	otherAuthor := other.(models.IAuthor)

	return aeq.author.IsEqual(otherAuthor)
}

// String _
func (aeq authorEq) String() string {
	return fmt.Sprintf("%#v", aeq.author)
}

func AuthorEq(author models.IAuthor) gomock.Matcher {
	return authorEq{author: author}
}

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

func remotifyDealer(localDealer models.IDealer) (*schema.ClientWithResponses, error) {
	e := handlers.NewDealerAPIRouter(context.Background(), localDealer, authoritySecret, sessionSecret)

	te := newEchoClientWrap(e)

	cl, err := schema.NewClientWithResponses("http://localhost:8080", schema.WithHTTPClient(te))

	if err != nil {
		return nil, err
	}

	return cl, nil
}

func createAuthor(t *testing.T) models.IAuthor {
	author := &models.Author{Name: "remote_dealer_test"}

	authorityToken, err := handlers.CreateAuthorityToken(authoritySecret, author.GetName())

	if err != nil {
		t.Fatalf("Failed to build author: %s", err)
		return nil
	}

	author.SetAuthorityKey(authorityToken)

	return author
}

func Test__AllocateSegment(t *testing.T) {
	author := createAuthor(t)

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
		AllocateSegment(gomock.Any(), AuthorEq(author), convertSegRequest).
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient)

		seg, err := rd.AllocateSegment(context.Background(), author, convertSegRequest)

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}

func Test__GetSegmentByID(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localDealer.
		EXPECT().
		GetSegmentByID(gomock.Any(), AuthorEq(author), "123").
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient)

		seg, err := rd.GetSegmentByID(context.Background(), author, "123")

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}
