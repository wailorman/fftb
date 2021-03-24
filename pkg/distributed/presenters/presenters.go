package presenters

import "github.com/wailorman/fftb/pkg/media/convert"

// Authority _
type Authority struct {
	Key string `json:"key"`
}

// ConvertSegment _
type ConvertSegment struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	State    string         `json:"state"`
	Params   convert.Params `json:"params"`
	Muxer    string         `json:"muxer"`
	Position int            `json:"position"`
}
