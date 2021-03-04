package models

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertOrder _
type ConvertOrder struct {
	Identity   string         `json:"identity"`
	Type       string         `json:"type"`
	State      string         `json:"state"`
	InFile     files.Filer    `json:"in_file"`
	OutFile    files.Filer    `json:"out_file"`
	Params     convert.Params `json:"params"`
	Publisher  IAuthor        `json:"publisher"`
	SegmentIDs []string       `json:"segment_i_ds"`
}

// OrderStateQueued _
const OrderStateQueued = "queued"

// OrderStateInProgress _
const OrderStateInProgress = "in_progress"

// OrderStateFinished _
const OrderStateFinished = "finished"

// GetID _
func (co *ConvertOrder) GetID() string {
	return co.Identity
}

// // SetID _
// func (co *ConvertOrder) SetID(id string) {
// 	co.Identity = id
// }

// GetType _
func (co *ConvertOrder) GetType() string {
	return ConvertV1Type
}

// // SetType _
// func (co *ConvertOrder) SetType(orderType string) {
// 	co.Type = orderType
// }

// GetPayload _
func (co *ConvertOrder) GetPayload() (string, error) {
	b, err := json.Marshal(co)

	return string(b), err
}

// GetState _
func (co *ConvertOrder) GetState() string {
	return co.State
}

// // SetState _
// func (co *ConvertOrder) SetState(state string) {
// 	co.State = state
// }

// // GetSegments _
// func (co *ConvertOrder) GetSegments() []ISegment {
// 	segments := make([]ISegment, 0)

// 	for _, task := range co.Segments {
// 		segments = append(segments, task)
// 	}

// 	return segments
// }

// // Failed _
// func (co *ConvertOrder) Failed(err error) {
// 	// TODO:
// 	panic(err)
// }

// GetPublisher _
func (co *ConvertOrder) GetPublisher() IAuthor {
	return co.Publisher
}

// // SetPublisher _
// func (co *ConvertOrder) SetPublisher(publisher IAuthor) {
// 	co.Publisher = publisher
// }

// MatchPublisher _
func (co *ConvertOrder) MatchPublisher(publisher IAuthor) bool {
	if co.Publisher == nil {
		return false
	}

	if publisher == nil {
		return false
	}

	return co.Publisher == publisher
}

// GetSegmentIDs _
func (co *ConvertOrder) GetSegmentIDs() []string {
	return co.SegmentIDs
}

// Validate _
func (co ConvertOrder) Validate() error {
	stateErr := validation.ValidateStruct(&co,
		validation.Field(&co.State,
			validation.Required,
			validation.In(
				OrderStateQueued,
				OrderStateInProgress,
				OrderStateFinished)))

	if stateErr != nil {
		return stateErr
	}

	return validation.ValidateStruct(&co,
		validation.Field(&co.Type, validation.Required, validation.In(ConvertV1Type)),
		validation.Field(&co.Identity, validation.Required),
		validation.Field(&co.InFile, validation.Required),
		validation.Field(&co.OutFile, validation.Required),
		validation.Field(&co.Publisher, validation.Required))
}
