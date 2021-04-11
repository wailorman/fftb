package models

import "github.com/wailorman/fftb/pkg/media/convert"

// RemoteAuthority _
type RemoteAuthority struct {
	Authority string `json:"authority"`
}

// RemoteConvertSegment _
type RemoteConvertSegment struct {
	ID       string         `json:"id"`
	OrderID  string         `json:"order_id"`
	Type     string         `json:"type"`
	State    string         `json:"state"`
	Params   convert.Params `json:"params"`
	Muxer    string         `json:"muxer"`
	Position int            `json:"position"`
}
