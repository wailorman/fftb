package models

import (
	"encoding/json"
	"time"

	"github.com/wailorman/fftb/pkg/media/convert"
)

// SegmentStatePrepared _
const SegmentStatePrepared = "prepared"

// SegmentStatePublished _
const SegmentStatePublished = "published"

// SegmentStateFinished _
const SegmentStateFinished = "finished"

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

// MatchPublisher _
func (ct *ConvertSegment) MatchPublisher(publisher IAuthor) bool {
	if ct.Publisher == nil {
		return false
	}

	if publisher == nil {
		return false
	}

	return ct.Publisher == publisher
}

// MatchPerformer _
func (ct *ConvertSegment) MatchPerformer(performer IAuthor) bool {
	if ct.LockedBy == nil {
		return false
	}

	if !ct.GetIsLocked() {
		return false
	}

	if performer == nil {
		return false
	}

	return ct.LockedBy == performer
}

// Lock _
func (ct *ConvertSegment) Lock(performer IAuthor) {
	lockedUntil := time.Now().Add(SegmentLockDuration)

	ct.LockedBy = performer
	ct.LockedUntil = &lockedUntil
}

// Unlock _
func (ct *ConvertSegment) Unlock() {
	ct.LockedBy = nil
	ct.LockedUntil = nil
}
