package converters

import (
	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// FromRPCSegment converts RPC segment to internal object
func FromRPCSegment(gSeg *pb.Segment) (models.ISegment, error) {
	switch gSeg.Type {
	case *pb.SegmentType_CONVERT_V1.Enum():
		if gSeg.ConvertParams == nil {
			return nil, errors.New("Missing convertParams")
		}

		return &models.ConvertSegment{
			Type:     models.ConvertV1Type,
			Identity: gSeg.Id,

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
		return nil, models.NewErrUnknownType(gSeg.Type.Enum().String())
	}
}

// ToRPCSegment converts internal segment to rpc format
func ToRPCSegment(mSeg models.ISegment) (*pb.Segment, error) {
	switch tmSeg := mSeg.(type) {
	case *models.ConvertSegment:
		return &pb.Segment{
			Type: *pb.SegmentType_CONVERT_V1.Enum(),
			Id:   tmSeg.Identity,
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

// ToRPCStorageClaimRequest generates RPC storage claim request
func ToRPCStorageClaimRequest(authorization string, req models.StorageClaimRequest) *pb.StorageClaimRequest {
	return &pb.StorageClaimRequest{
		Authorization: authorization,
		SegmentId:     req.SegmentID,
		Purpose:       ToRPCStorageClaimPurpose(req.Purpose),
		Name:          req.Name,
	}
}

func FromRPCStorageClaimPurpose(purpose pb.StorageClaimPurpose) models.StorageClaimPurpose {
	switch purpose {
	case pb.StorageClaimPurpose_CONVERT_INPUT:
		return models.ConvertInputStorageClaimPurpose
	case pb.StorageClaimPurpose_CONVERT_OUTPUT:
		return models.ConvertOutputStorageClaimPurpose
	default:
		return models.NoneStorageClaimPurpose
	}
}

func ToRPCStorageClaimPurpose(purpose models.StorageClaimPurpose) pb.StorageClaimPurpose {
	switch purpose {
	case models.ConvertInputStorageClaimPurpose:
		return pb.StorageClaimPurpose_CONVERT_INPUT
	case models.ConvertOutputStorageClaimPurpose:
		return pb.StorageClaimPurpose_CONVERT_OUTPUT
	default:
		return pb.StorageClaimPurpose_NONE
	}
}

// ToRPCProgress converts internal progress message to rpc format
func ToRPCProgress(authorization, segmentID string, progress models.IProgress) (*pb.ProgressNotification, error) {
	switch progress.Step() {
	case models.UploadingInputStep:
		return &pb.ProgressNotification{
			Authorization: authorization,
			Step:          pb.ProgressNotification_UPLOADING_INPUT,
			Progress:      progress.Percent(),
			SegmentId:     segmentID}, nil

	case models.DownloadingInputStep:
		return &pb.ProgressNotification{
			Authorization: authorization,
			Step:          pb.ProgressNotification_DOWNLOADING_INPUT,
			Progress:      progress.Percent(),
			SegmentId:     segmentID}, nil

	case models.ProcessingStep:
		return &pb.ProgressNotification{
			Authorization: authorization,
			Step:          pb.ProgressNotification_PROCESSING,
			Progress:      progress.Percent(),
			SegmentId:     segmentID}, nil

	case models.UploadingOutputStep:
		return &pb.ProgressNotification{
			Authorization: authorization,
			Step:          pb.ProgressNotification_UPLOADING_OUTPUT,
			Progress:      progress.Percent(),
			SegmentId:     segmentID}, nil

	case models.DownloadingOutputStep:
		return &pb.ProgressNotification{
			Authorization: authorization,
			Step:          pb.ProgressNotification_DOWNLOADING_OUTPUT,
			Progress:      progress.Percent(),
			SegmentId:     segmentID}, nil

	default:
		return nil, models.NewErrUnknownType(string(progress.Step()))
	}
}
