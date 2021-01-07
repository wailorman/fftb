package models

import (
	"encoding/json"
	"time"

	"github.com/wailorman/fftb/pkg/media/convert"
)

// SegmentPreparedState _
const SegmentPreparedState = "prepared"

// SegmentPublishedState _
const SegmentPublishedState = "published"

// SegmentFinishedState _
const SegmentFinishedState = "finished"

// ConvertSegment _
type ConvertSegment struct {
	Identity                   string
	OrderIdentity              string
	Type                       string
	InputStorageClaimIdentity  string
	OutputStorageClaimIdentity string
	State                      string

	Params      convert.Params
	LockedUntil *time.Time
	LockedBy    string
	Muxer       string
	// VideoCodec string
	// // HWAccel          string
	// // VideoBitRate     string
	// VideoQuality int
	// // Preset           string
	// // Scale            string
	// // KeyframeInterval int
}

// GetID _
func (ct *ConvertSegment) GetID() string {
	return ct.Identity
}

// GetType _
func (ct *ConvertSegment) GetType() string {
	return ConvertV1Type
}

// GetOrderID _
func (ct *ConvertSegment) GetOrderID() string {
	return ct.OrderIdentity
}

// GetInputStorageClaimIdentity _
func (ct *ConvertSegment) GetInputStorageClaimIdentity() string {
	return ct.InputStorageClaimIdentity
}

// GetOutputStorageClaimIdentity _
func (ct *ConvertSegment) GetOutputStorageClaimIdentity() string {
	return ct.OutputStorageClaimIdentity
}

// GetPayload _
func (ct *ConvertSegment) GetPayload() (string, error) {
	b, err := json.Marshal(ct)

	return string(b), err
}

// GetIsLocked _
func (ct *ConvertSegment) GetIsLocked() bool {
	if ct.LockedUntil == nil {
		return false
	}

	return time.Now().After(*ct.LockedUntil) && ct.LockedBy != ""
}

// GetLockedBy _
func (ct *ConvertSegment) GetLockedBy() string {
	if !ct.GetIsLocked() {
		return ""
	}

	return ct.LockedBy
}

// // GetStorageClaim _
// func (ct *ConvertSegment) GetStorageClaim() IStorageClaim {
// 	return ct.StorageClaim
// }

// Failed _
func (ct *ConvertSegment) Failed(err error) {
	// TODO:
	// panic(ErrNotImplemented)
	panic(err)
	// return
}

// GetState _
func (ct *ConvertSegment) GetState() string {
	return ct.State
}
