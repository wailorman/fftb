package remote_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/local"
	"github.com/wailorman/fftb/pkg/distributed/models"
	mock_models "github.com/wailorman/fftb/pkg/distributed/models/mocks"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	dealerSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
	"github.com/wailorman/fftb/pkg/media/convert"
)

func remotifyDealerApi(t *testing.T, localDealer models.IDealer) *dealerSchema.ClientWithResponses {
	e := handlers.NewDealerAPIRouter(context.Background(), logrus.New(), localDealer, authoritySecret, sessionSecret)

	te := newEchoClientWrap(e)

	cl, err := dealerSchema.NewClientWithResponses("http://localhost:8080", dealerSchema.WithHTTPClient(te))

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return cl
}

func Test__AllocateSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)
	localStorageClient := mock_models.NewMockIStorageClient(ctrl)

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

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	seg, err := rd.AllocateSegment(context.Background(), author, convertSegRequest)

	if assert.NoError(t, err) {
		assert.Equal(t, "123", seg.GetID())
	}
}

func Test__GetSegmentByID(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localStorageClient := mock_models.NewMockIStorageClient(ctrl)

	localDealer.
		EXPECT().
		GetSegmentByID(gomock.Any(), AuthorEq(author), "123").
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	seg, err := rd.GetSegmentByID(context.Background(), author, "123")

	if assert.NoError(t, err) {
		assert.Equal(t, "123", seg.GetID())
	}
}

func Test__FindFreeSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localStorageClient := mock_models.NewMockIStorageClient(ctrl)

	localDealer.
		EXPECT().
		FindFreeSegment(gomock.Any(), AuthorEq(author)).
		Return(&models.ConvertSegment{Identity: "123"}, nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	seg, err := rd.FindFreeSegment(context.Background(), author)

	if assert.NoError(t, err) {
		assert.Equal(t, "123", seg.GetID())
	}
}

func Test__FailSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localStorageClient := mock_models.NewMockIStorageClient(ctrl)
	reportedErr := errors.New("Testing error")
	segmentID := "123"

	localDealer.
		EXPECT().
		FailSegment(gomock.Any(), AuthorEq(author), segmentID, HasString(reportedErr.Error())).
		Return(nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	err := rd.FailSegment(context.Background(), author, segmentID, reportedErr)

	assert.NoError(t, err)
}

func Test__FinishSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localStorageClient := mock_models.NewMockIStorageClient(ctrl)
	segmentID := "123"

	localDealer.
		EXPECT().
		FinishSegment(gomock.Any(), AuthorEq(author), segmentID).
		Return(nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	err := rd.FinishSegment(context.Background(), author, segmentID)

	assert.NoError(t, err)
}

func Test__QuitSegment(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	localStorageClient := mock_models.NewMockIStorageClient(ctrl)
	segmentID := "123"

	localDealer.
		EXPECT().
		QuitSegment(gomock.Any(), AuthorEq(author), segmentID).
		Return(nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, localStorageClient, authoritySecret)

	err := rd.QuitSegment(context.Background(), author, segmentID)

	assert.NoError(t, err)
}

func Test__GetInputStorageClaim(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	segmentID := "123"
	storageClient := local.NewStorageClient(".")
	storageClaim, err := storageClient.BuildStorageClaimByURL("file://remote_dealer_test.go")

	if err != nil {
		t.Errorf("Failed to build storage claim: %s", err)
	}

	localDealer.
		EXPECT().
		GetInputStorageClaim(gomock.Any(), AuthorEq(author), segmentID).
		Return(storageClaim, nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, storageClient, authoritySecret)

	resStorageClaim, err := rd.GetInputStorageClaim(context.Background(), author, segmentID)

	if assert.NoError(t, err) {
		assert.Equal(t, storageClaim.GetID(), resStorageClaim.GetID())
		assert.NotEqual(t, "", resStorageClaim.GetID())
	}
}

func Test__AllocateOutputStorageClaim(t *testing.T) {
	author := createAuthor(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localDealer := mock_models.NewMockIDealer(ctrl)

	segmentID := "123"
	storageController := local.NewStorageClient(".")
	storageClaim, err := storageController.BuildStorageClaimByURL("file://remote_dealer_test.go")

	if err != nil {
		t.Errorf("Failed to build storage claim: %s", err)
	}

	localDealer.
		EXPECT().
		AllocateOutputStorageClaim(gomock.Any(), AuthorEq(author), segmentID).
		Return(storageClaim, nil)

	apiClient := remotifyDealerApi(t, localDealer)

	rd := remote.NewDealer(apiClient, storageController, authoritySecret)

	resStorageClaim, err := rd.AllocateOutputStorageClaim(context.Background(), author, segmentID)

	if assert.NoError(t, err) {
		assert.Equal(t, storageClaim.GetID(), resStorageClaim.GetID())
		assert.NotEqual(t, "", resStorageClaim.GetID())
	}
}
