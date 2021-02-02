package models

import (
	"encoding/json"

	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertOrder _
type ConvertOrder struct {
	Identity string
	Type     string
	State    string
	// Segments []*ConvertSegment

	Params convert.Params

	Publisher IAuthor
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
