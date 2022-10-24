package remote_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/models"
	mock_models "github.com/wailorman/fftb/pkg/distributed/models/mocks"
	"github.com/wailorman/fftb/pkg/distributed/remote"
	"github.com/wailorman/fftb/pkg/distributed/remote/handlers"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/distributed/test/factories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var f *factories.Builder
var ctrl *gomock.Controller
var conn *grpc.ClientConn
var grpcServer *grpc.Server
var localDealer *mock_models.MockIDealer
var remotedDealer models.IDealer
var storageClient *mock_models.MockIStorageClient

var progressPercent = 0.34

func makeBufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
func remotifyDealer(t *testing.T, localDealer models.IDealer, storageClient models.IStorageClient) models.IDealer {
	bufSize := 1024 * 1024 * 10
	lis := bufconn.Listen(bufSize)
	grpcServer = grpc.NewServer()
	pb.RegisterDealerServer(grpcServer, handlers.NewDealerHandler(localDealer, nil, nil))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(makeBufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := pb.NewDealerClient(conn)
	return remote.NewDealer(client, storageClient)
}

func setup(t *testing.T) {
	f = factories.NewBuilder()
	ctrl = gomock.NewController(t)

	localDealer = mock_models.NewMockIDealer(ctrl)
	storageClient = mock_models.NewMockIStorageClient(ctrl)
	remotedDealer = remotifyDealer(t, localDealer, storageClient)
}

func teardown() {
	if conn != nil {
		conn.Close()
	}

	if grpcServer != nil {
		grpcServer.Stop()
	}

	ctrl.Finish()
}

func Test__AllocateSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			AllocateSegment(gomock.Any(), f.Author, f.ConvertDealerRequest()).
			Return(f.ConvertSegment(), nil)

		actualSegment, err :=
			remotedDealer.AllocateSegment(context.Background(), f.Author, f.ConvertDealerRequest())

		if assert.NoError(t, err) {
			assert.Equal(t, f.ConvertSegment(), actualSegment)
		}
	})
}

func Test__GetOutputStorageClaim(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			GetOutputStorageClaim(gomock.Any(), f.Author, f.SegmentID).
			Return(f.StorageClaim(), nil)

		storageClient.
			EXPECT().
			BuildStorageClaimByURL(f.StorageClaimURL).
			Return(f.StorageClaim(), nil)

		actualStorageClaim, err :=
			remotedDealer.GetOutputStorageClaim(context.Background(), f.Author, f.SegmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, f.StorageClaim(), actualStorageClaim)
		}
	})
}
func Test__AllocateInputStorageClaim(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			AllocateInputStorageClaim(gomock.Any(), f.Author, f.SegmentID).
			Return(f.StorageClaim(), nil)

		storageClient.
			EXPECT().
			BuildStorageClaimByURL(f.StorageClaimURL).
			Return(f.StorageClaim(), nil)

		actualStorageClaim, err :=
			remotedDealer.AllocateInputStorageClaim(context.Background(), f.Author, f.SegmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, f.StorageClaim(), actualStorageClaim)
		}
	})
}
func Test__GetInputStorageClaim(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			GetInputStorageClaim(gomock.Any(), f.Author, f.SegmentID).
			Return(f.StorageClaim(), nil)

		storageClient.
			EXPECT().
			BuildStorageClaimByURL(f.StorageClaimURL).
			Return(f.StorageClaim(), nil)

		actualStorageClaim, err :=
			remotedDealer.GetInputStorageClaim(context.Background(), f.Author, f.SegmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, f.StorageClaim(), actualStorageClaim)
		}
	})
}
func Test__AllocateOutputStorageClaim(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			AllocateOutputStorageClaim(gomock.Any(), f.Author, f.SegmentID).
			Return(f.StorageClaim(), nil)

		storageClient.
			EXPECT().
			BuildStorageClaimByURL(f.StorageClaimURL).
			Return(f.StorageClaim(), nil)

		actualStorageClaim, err :=
			remotedDealer.AllocateOutputStorageClaim(context.Background(), f.Author, f.SegmentID)

		if assert.NoError(t, err) {
			assert.Equal(t, f.StorageClaim(), actualStorageClaim)
		}
	})
}

func Test__NotifyRawUpload(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		mProgress := f.Progress(models.UploadingInputStep, progressPercent)

		localDealer.
			EXPECT().
			NotifyRawUpload(gomock.Any(), f.Author, f.SegmentID, mProgress).
			Return(nil)

		err := remotedDealer.NotifyRawUpload(context.Background(), f.Author, f.SegmentID, mProgress)

		assert.NoError(t, err)
	})
}
func Test__NotifyRawDownload(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		mProgress := f.Progress(models.DownloadingInputStep, progressPercent)

		localDealer.
			EXPECT().
			NotifyRawDownload(gomock.Any(), f.Author, f.SegmentID, mProgress).
			Return(nil)

		err := remotedDealer.NotifyRawDownload(context.Background(), f.Author, f.SegmentID, mProgress)

		assert.NoError(t, err)
	})
}
func Test__NotifyProcess(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		mProgress := f.Progress(models.ProcessingStep, progressPercent)

		localDealer.
			EXPECT().
			NotifyProcess(gomock.Any(), f.Author, f.SegmentID, mProgress).
			Return(nil)

		err := remotedDealer.NotifyProcess(context.Background(), f.Author, f.SegmentID, mProgress)

		assert.NoError(t, err)
	})
}
func Test__NotifyResultUpload(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		mProgress := f.Progress(models.UploadingOutputStep, progressPercent)

		localDealer.
			EXPECT().
			NotifyResultUpload(gomock.Any(), f.Author, f.SegmentID, mProgress).
			Return(nil)

		err := remotedDealer.NotifyResultUpload(context.Background(), f.Author, f.SegmentID, mProgress)

		assert.NoError(t, err)
	})
}
func Test__NotifyResultDownload(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		mProgress := f.Progress(models.DownloadingOutputStep, progressPercent)

		localDealer.
			EXPECT().
			NotifyResultDownload(gomock.Any(), f.Author, f.SegmentID, mProgress).
			Return(nil)

		err := remotedDealer.NotifyResultDownload(context.Background(), f.Author, f.SegmentID, mProgress)

		assert.NoError(t, err)
	})
}

func Test__PublishSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			PublishSegment(gomock.Any(), f.Author, f.SegmentID).
			Return(nil)

		err := remotedDealer.PublishSegment(context.Background(), f.Author, f.SegmentID)

		assert.NoError(t, err)
	})
}

func Test__RepublishSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			RepublishSegment(gomock.Any(), f.Author, f.SegmentID).
			Return(nil)

		err := remotedDealer.RepublishSegment(context.Background(), f.Author, f.SegmentID)

		assert.NoError(t, err)
	})
}

func Test__AcceptSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			AcceptSegment(gomock.Any(), f.Author, f.SegmentID).
			Return(nil)

		err := remotedDealer.AcceptSegment(context.Background(), f.Author, f.SegmentID)

		assert.NoError(t, err)
	})
}

func Test__FinishSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			FinishSegment(gomock.Any(), f.Author, f.SegmentID).
			Return(nil)

		err := remotedDealer.FinishSegment(context.Background(), f.Author, f.SegmentID)

		assert.NoError(t, err)
	})
}

func Test__QuitSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			QuitSegment(gomock.Any(), f.Author, f.SegmentID).
			Return(nil)

		err := remotedDealer.QuitSegment(context.Background(), f.Author, f.SegmentID)

		assert.NoError(t, err)
	})
}

func Test__CancelSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			CancelSegment(gomock.Any(), f.Author, f.SegmentID, f.CancellationReason).
			Return(nil)

		err := remotedDealer.CancelSegment(context.Background(), f.Author, f.SegmentID, f.CancellationReason)

		assert.NoError(t, err)
	})
}

func Test__FailSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			FailSegment(gomock.Any(), f.Author, f.SegmentID, gomock.Any()).
			Return(nil)
		err := remotedDealer.FailSegment(context.Background(), f.Author, f.SegmentID, f.FailureReason)

		assert.NoError(t, err)
	})
}

func Test__FindFreeSegment(t *testing.T) {
	t.Run("calls local dealer", func(t *testing.T) {
		setup(t)
		defer teardown()

		localDealer.
			EXPECT().
			FindFreeSegment(gomock.Any(), f.Author).
			Return(f.ConvertSegment(), nil)

		actualSegment, err := remotedDealer.FindFreeSegment(context.Background(), f.Author)

		if assert.NoError(t, err) {
			assert.Equal(t, f.ConvertSegment(), actualSegment)
		}
	})
}
