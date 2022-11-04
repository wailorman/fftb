package models

import (
	"github.com/wailorman/fftb/pkg/media/convert"
)

// SegmentStatePrepared _
const SegmentStatePrepared = "prepared"

// SegmentStatePublished _
const SegmentStatePublished = "published"

// SegmentStateInProgress is dynamic state, used only in some presenters.
// Can be returned by GetCurrentState()
const SegmentStateInProgress = "in_progress"

// SegmentStateAccepted _
const SegmentStateAccepted = "accepted"

// SegmentStateFinished _
const SegmentStateFinished = "finished"

// SegmentStateCancelled _
const SegmentStateCancelled = "cancelled"

// ConvertSegment _
type ConvertSegment struct {
	Identity string `json:"identity"`
	Type     string `json:"type"`

	Params   convert.Params `json:"params"`
	Muxer    string         `json:"muxer"`
	Position int            `json:"position"`
}

// GetID _
func (ct *ConvertSegment) GetID() string {
	return ct.Identity
}

// GetType _
func (ct *ConvertSegment) GetType() string {
	return ConvertV1Type
}
