package converters

import (
	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// FromRPCDealerRequest converts RPC dealer request to internal abstract format. Returns authorization string, dealer request & error
func FromRPCDealerRequest(gReq *pb.DealerRequest) (string, models.IDealerRequest, error) {
	switch gReq.Type {
	case models.ConvertV1Type:
		if gReq.ConvertParams == nil {
			return "", nil, errors.New("Missing .convertParams")
		}

		return gReq.Authorization, &models.ConvertDealerRequest{
			Type:          models.ConvertV1Type,
			Identity:      gReq.Id,
			OrderIdentity: gReq.OrderId,

			Params: convert.Params{
				VideoCodec:       gReq.ConvertParams.VideoCodec,
				HWAccel:          gReq.ConvertParams.HwAccel,
				VideoBitRate:     gReq.ConvertParams.VideoBitRate,
				VideoQuality:     int(gReq.ConvertParams.VideoQuality),
				Preset:           gReq.ConvertParams.Preset,
				Scale:            gReq.ConvertParams.Scale,
				KeyframeInterval: int(gReq.ConvertParams.KeyframeInterval),
			},

			Muxer:    gReq.ConvertParams.Muxer,
			Position: int(gReq.ConvertParams.Position),
		}, nil

	default:
		return "", nil, models.NewErrUnknownType(gReq.Type)
	}
}

// ToRPCDealerRequest converts internal dealer request to RPC format
func ToRPCDealerRequest(authorization string, mReq models.IDealerRequest) (*pb.DealerRequest, error) {
	switch tmReq := mReq.(type) {
	case *models.ConvertDealerRequest:
		return &pb.DealerRequest{
			Authorization: authorization,
			Type:          models.ConvertV1Type,
			Id:            tmReq.Identity,
			OrderId:       tmReq.OrderIdentity,
			ConvertParams: &pb.ConvertSegmentParams{
				VideoCodec:       tmReq.Params.VideoCodec,
				HwAccel:          tmReq.Params.HWAccel,
				VideoBitRate:     tmReq.Params.VideoBitRate,
				VideoQuality:     int32(tmReq.Params.VideoQuality),
				Preset:           tmReq.Params.Preset,
				Scale:            tmReq.Params.Scale,
				KeyframeInterval: int32(tmReq.Params.KeyframeInterval),
				Muxer:            tmReq.Muxer,
				Position:         int32(tmReq.Position),
			},
		}, nil

	default:
		return nil, models.NewErrUnknownType(mReq.GetType())
	}
}

// FromRPCSegment converts RPC segment to internal object
func FromRPCSegment(gSeg *pb.Segment) (models.ISegment, error) {
	switch gSeg.Type {
	case models.ConvertV1Type:
		if gSeg.ConvertParams == nil {
			return nil, errors.New("Missing convertParams")
		}

		return &models.ConvertSegment{
			Type:          models.ConvertV1Type,
			Identity:      gSeg.Id,
			OrderIdentity: gSeg.OrderId,

			Params: convert.Params{
				VideoCodec:       gSeg.ConvertParams.VideoCodec,
				HWAccel:          gSeg.ConvertParams.HwAccel,
				VideoBitRate:     gSeg.ConvertParams.VideoBitRate,
				VideoQuality:     int(gSeg.ConvertParams.VideoQuality),
				Preset:           gSeg.ConvertParams.Preset,
				Scale:            gSeg.ConvertParams.Scale,
				KeyframeInterval: int(gSeg.ConvertParams.KeyframeInterval),
			},

			Muxer:    gSeg.ConvertParams.Muxer,
			Position: int(gSeg.ConvertParams.Position),
		}, nil

	default:
		return nil, models.NewErrUnknownType(gSeg.Type)
	}
}

// ToRPCSegment converts internal segment to rpc format
func ToRPCSegment(mSeg models.ISegment) (*pb.Segment, error) {
	switch tmSeg := mSeg.(type) {
	case *models.ConvertSegment:
		return &pb.Segment{
			Type:    models.ConvertV1Type,
			Id:      tmSeg.Identity,
			OrderId: tmSeg.OrderIdentity,
			ConvertParams: &pb.ConvertSegmentParams{
				VideoCodec:       tmSeg.Params.VideoCodec,
				HwAccel:          tmSeg.Params.HWAccel,
				VideoBitRate:     tmSeg.Params.VideoBitRate,
				VideoQuality:     int32(tmSeg.Params.VideoQuality),
				Preset:           tmSeg.Params.Preset,
				Scale:            tmSeg.Params.Scale,
				KeyframeInterval: int32(tmSeg.Params.KeyframeInterval),
				Muxer:            tmSeg.Muxer,
				Position:         int32(tmSeg.Position),
			},
		}, nil
	default:
		return nil, models.NewErrUnknownType(mSeg.GetType())
	}
}

// ToRPCStorageClaim converts internal storage claim to rpc format
func ToRPCStorageClaim(sc models.IStorageClaim) *pb.StorageClaim {
	return &pb.StorageClaim{
		Url: sc.GetURL(),
	}
}

// ToRPCStorageClaimRequest generates RPC storage claim request
func ToRPCStorageClaimRequest(authorization string, segmentID string) *pb.StorageClaimRequest {
	return &pb.StorageClaimRequest{
		Authorization: authorization,
		SegmentId:     segmentID,
	}
}
