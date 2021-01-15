package models

import (
	"encoding/json"

	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertOrder _
type ConvertOrder struct {
	Identity string
	Type     string
	Segments []*ConvertSegment

	Params    convert.Params
	Publisher IAuthor
}

// GetID _
func (co *ConvertOrder) GetID() string {
	return co.Identity
}

// GetType _
func (co *ConvertOrder) GetType() string {
	return ConvertV1Type
}

// GetPayload _
func (co *ConvertOrder) GetPayload() (string, error) {
	b, err := json.Marshal(co)

	return string(b), err
}

// GetSegments _
func (co *ConvertOrder) GetSegments() []ISegment {
	segments := make([]ISegment, 0)

	for _, task := range co.Segments {
		segments = append(segments, task)
	}

	return segments
}

// // Failed _
// func (co *ConvertOrder) Failed(err error) {
// 	// TODO:
// 	panic(err)
// }

// GetPublisher _
func (co *ConvertOrder) GetPublisher() IAuthor {
	return co.Publisher
}
