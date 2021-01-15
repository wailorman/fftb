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

	Params convert.Params
	Muxer  string

	Publisher   IAuthor
	LockedUntil *time.Time
	LockedBy    IAuthor
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
	if ct.LockedUntil == nil || ct.LockedBy == nil {
		return false
	}

	return time.Now().After(*ct.LockedUntil)
}

// GetLockedBy _
func (ct *ConvertSegment) GetLockedBy() IAuthor {
	if !ct.GetIsLocked() {
		return nil
	}

	return ct.LockedBy
}

// GetLockedUntil _
func (ct *ConvertSegment) GetLockedUntil() *time.Time {
	if !ct.GetIsLocked() {
		return nil
	}

	return ct.LockedUntil
}

// // GetStorageClaim _
// func (ct *ConvertSegment) GetStorageClaim() IStorageClaim {
// 	return ct.StorageClaim
// }

// Failed _
// func (ct *ConvertSegment) Failed(err error) {
// 	// TODO:
// 	// panic(ErrNotImplemented)
// 	panic(err)
// 	// return
// }

// GetState _
func (ct *ConvertSegment) GetState() string {
	return ct.State
}

// GetPublisher _
func (ct *ConvertSegment) GetPublisher() IAuthor {
	return ct.Publisher
}

// GetPerformer _
func (ct *ConvertSegment) GetPerformer() IAuthor {
	return ct.LockedBy
}
