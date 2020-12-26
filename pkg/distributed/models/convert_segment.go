package models

import (
	"encoding/json"

	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertSegment _
type ConvertSegment struct {
	Identity             string
	OrderIdentity        string
	Type                 string
	StorageClaimIdentity string

	Params convert.ConverterTask
	// Muxer      string
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

// GetStorageClaimIdentity _
func (ct *ConvertSegment) GetStorageClaimIdentity() string {
	return ct.StorageClaimIdentity
}

// GetPayload _
func (ct *ConvertSegment) GetPayload() (string, error) {
	b, err := json.Marshal(ct)

	return string(b), err
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
