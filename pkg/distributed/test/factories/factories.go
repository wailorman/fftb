package factories

import (
	"errors"

	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/distributed/s3"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// Builder _
type Builder struct {
	Authorization string

	SegmentID string
	OrderID   string

	VideoCodec       string
	HWAccel          string
	VideoBitRate     string
	VideoQuality     int
	Preset           string
	Scale            string
	KeyframeInterval int

	Muxer    string
	Position int

	Author models.IAuthor

	StorageClaimID   string
	StorageClaimURL  string
	StorageClaimSize int

	CancellationReason string
	FailureReason      error
}

// NewBuilder _
func NewBuilder() *Builder {
	return &Builder{
		Authorization: "SD8F79ASD87F6SA9D8F7SA9D8F7",

		SegmentID: "b9ef0c17-c69d-4666-9461-287be9e1657a",
		OrderID:   "f6ffc440-b1f6-4d6c-8911-a5f879bd39e9",

		VideoCodec:       "h264",
		HWAccel:          "nvenc",
		VideoBitRate:     "10M",
		VideoQuality:     30,
		Preset:           "slow",
		Scale:            "1/2",
		KeyframeInterval: 30,

		Muxer:    "mp4",
		Position: 49,

		Author: models.LocalAuthor,

		StorageClaimID:   "s3://bucket/test.mp4",
		StorageClaimURL:  "https://s3.example.com/bucket/test.mp4",
		StorageClaimSize: 9999,

		CancellationReason: "Just something went wrong",
		FailureReason:      errors.New("Something failed"),
	}
}

// ConvertParams _
func (b *Builder) ConvertParams() convert.Params {
	return convert.Params{
		VideoCodec:       b.VideoCodec,
		HWAccel:          b.HWAccel,
		VideoBitRate:     b.VideoBitRate,
		VideoQuality:     b.VideoQuality,
		Preset:           b.Preset,
		Scale:            b.Scale,
		KeyframeInterval: b.KeyframeInterval,
	}
}

// RPCConvertParams _
func (b *Builder) RPCConvertParams() *pb.ConvertSegmentParams {
	return &pb.ConvertSegmentParams{
		VideoCodec:       b.VideoCodec,
		HwAccel:          b.HWAccel,
		VideoBitRate:     b.VideoBitRate,
		VideoQuality:     int32(b.VideoQuality),
		Preset:           b.Preset,
		Scale:            b.Scale,
		KeyframeInterval: int32(b.KeyframeInterval),
		Muxer:            b.Muxer,
		Position:         int32(b.Position),
	}
}

// ConvertDealerRequest _
func (b *Builder) ConvertDealerRequest() models.IDealerRequest {
	return &models.ConvertDealerRequest{
		Type:          models.ConvertV1Type,
		Identity:      b.SegmentID,
		OrderIdentity: b.OrderID,
		Params:        b.ConvertParams(),
		Muxer:         b.Muxer,
		Position:      b.Position,
	}
}

// RPCConvertDealerRequest _
func (b *Builder) RPCConvertDealerRequest() *pb.DealerRequest {
	return &pb.DealerRequest{
		Authorization: b.Authorization,
		Type:          models.ConvertV1Type,
		Id:            b.SegmentID,
		OrderId:       b.OrderID,
		ConvertParams: b.RPCConvertParams(),
	}
}

// ConvertSegment _
func (b *Builder) ConvertSegment() models.ISegment {
	return &models.ConvertSegment{
		Identity:      b.SegmentID,
		OrderIdentity: b.OrderID,
		Type:          models.ConvertV1Type,
		Params:        b.ConvertParams(),
		Muxer:         b.Muxer,
		Position:      b.Position,
	}
}

// RPCConvertSegment _
func (b *Builder) RPCConvertSegment() *pb.Segment {
	return &pb.Segment{
		Type:          models.ConvertV1Type,
		Id:            b.SegmentID,
		OrderId:       b.OrderID,
		ConvertParams: b.RPCConvertParams(),
	}
}

// StorageClaim _
func (b *Builder) StorageClaim() models.IStorageClaim {
	return s3.BuildStorageClaim(b.StorageClaimID, b.StorageClaimURL, b.StorageClaimSize)
}

// RPCStorageClaim _
func (b *Builder) RPCStorageClaim() *pb.StorageClaim {
	return &pb.StorageClaim{Url: b.StorageClaimURL}
}

// Progress _
func (b *Builder) Progress(step models.ProgressStep, progress float64) models.IProgress {
	return dlog.BuildProgress(step, progress)
}

// RPCProgress _
func (b *Builder) RPCProgress(step pb.ProgressNotification_Step, progress float64) *pb.ProgressNotification {
	return &pb.ProgressNotification{
		Authorization: b.Authorization,
		Step:          step,
		Progress:      progress,
		SegmentId:     b.SegmentID,
	}
}
