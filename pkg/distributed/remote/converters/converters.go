package converters

import (
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

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
