package remote_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	mock_models "github.com/wailorman/fftb/pkg/distributed/models/mocks"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/distributed/schema"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

var authoritySecret = []byte("authority_secret_remote_dealer")
var sessionSecret = []byte("session_secret_remote_dealer")

type authorEq struct {
	author models.IAuthor
}

// Matches _
func (aeq authorEq) Matches(other interface{}) bool {
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

// hasString _
type hasString struct {
	subStr string
}

func HasString(subStr string) gomock.Matcher {
	return &hasString{subStr: subStr}
}

// String _
func (hs hasString) String() string {
	return fmt.Sprintf("Contains %s", hs.subStr)
}

// Matches _
func (hs hasString) Matches(other interface{}) bool {
	var otherStr string

	if otherStrStringer, ok := other.(fmt.Stringer); ok {
		otherStr = otherStrStringer.String()
	}

	if otherStrStringer, ok := other.(error); ok {
		otherStr = otherStrStringer.Error()
	}

	return strings.Contains(otherStr, hs.subStr)
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
	e := handlers.NewDealerAPIRouter(context.Background(), logrus.New(), localDealer, authoritySecret, sessionSecret)

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
	localStorageControl := mock_models.NewMockIStorageController(ctrl)

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
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

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
	localStorageControl := mock_models.NewMockIStorageController(ctrl)

	localDealer.
		EXPECT().
		GetSegmentByID(gomock.Any(), AuthorEq(author), "123").
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

		seg, err := rd.GetSegmentByID(context.Background(), author, "123")

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}

func Test__FindFreeSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	localStorageControl := mock_models.NewMockIStorageController(ctrl)

	localDealer.
		EXPECT().
		FindFreeSegment(gomock.Any(), AuthorEq(author)).
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

		seg, err := rd.FindFreeSegment(context.Background(), author)

		if assert.NoError(t, err) {
			assert.Equal(t, "123", seg.GetID())
		}
	}
}

func Test__FailSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	localStorageControl := mock_models.NewMockIStorageController(ctrl)
	reportedErr := errors.New("Testing error")
	segmentID := "123"

	localDealer.
		EXPECT().
		FailSegment(gomock.Any(), AuthorEq(author), segmentID, HasString(reportedErr.Error())).
		Return(nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

		err := rd.FailSegment(context.Background(), author, segmentID, reportedErr)

		assert.NoError(t, err)
	}
}

func Test__FinishSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	localStorageControl := mock_models.NewMockIStorageController(ctrl)
	segmentID := "123"

	localDealer.
		EXPECT().
		FinishSegment(gomock.Any(), AuthorEq(author), segmentID).
		Return(nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

		err := rd.FinishSegment(context.Background(), author, segmentID)

		assert.NoError(t, err)
	}
}

func Test__QuitSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	localStorageControl := mock_models.NewMockIStorageController(ctrl)
	segmentID := "123"

	localDealer.
		EXPECT().
		QuitSegment(gomock.Any(), AuthorEq(author), segmentID).
		Return(nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, localStorageControl, authoritySecret)

		err := rd.QuitSegment(context.Background(), author, segmentID)

		assert.NoError(t, err)
	}
}

func Test__GetInputStorageClaim(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	segmentID := "123"
	storageController := local.NewStorageControl(files.NewPath("."))
	storageClaim, err := storageController.BuildStorageClaim("remote_dealer_test.go")

	if err != nil {
		t.Errorf("Failed to build storage claim: %s", err)
	}

	localDealer.
		EXPECT().
		GetInputStorageClaim(gomock.Any(), AuthorEq(author), segmentID).
		Return(storageClaim, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, storageController, authoritySecret)

		resStorageClaim, err := rd.GetInputStorageClaim(context.Background(), author, segmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, storageClaim.GetID(), resStorageClaim.GetID())
			assert.NotEqual(t, "", resStorageClaim.GetID())
		}
	}
}

func Test__AllocateOutputStorageClaim(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	segmentID := "123"
	storageController := local.NewStorageControl(files.NewPath("."))
	storageClaim, err := storageController.BuildStorageClaim("remote_dealer_test.go")

	if err != nil {
		t.Errorf("Failed to build storage claim: %s", err)
	}

	localDealer.
		EXPECT().
		AllocateOutputStorageClaim(gomock.Any(), AuthorEq(author), segmentID).
		Return(storageClaim, nil)

	apiClient, err := remotifyDealer(localDealer)

	if assert.NoError(t, err) {
		rd := remote.NewDealer(apiClient, storageController, authoritySecret)

		resStorageClaim, err := rd.AllocateOutputStorageClaim(context.Background(), author, segmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, storageClaim.GetID(), resStorageClaim.GetID())
			assert.NotEqual(t, "", resStorageClaim.GetID())
		}
	}
}
