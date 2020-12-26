package models

import (
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertV1Type _
const ConvertV1Type = "convert/v1"

// ConvertContracterRequest _
type ConvertContracterRequest struct {
	Params convert.ConverterTask
	// Muxer      string
	// VideoCodec string
	// // HWAccel          string
	// // VideoBitRate     string
	// VideoQuality int
	// // Preset           string
	// // Scale            string
	// // KeyframeInterval int
	// InFile files.Filer
}

// GetType _
func (cr *ConvertContracterRequest) GetType() string {
	return ConvertV1Type
}

// ConvertDealerRequest _
type ConvertDealerRequest struct {
	Type          string
	Identity      string
	OrderIdentity string

	Params convert.ConverterTask
	// Muxer      string
	// VideoCodec string
	// // HWAccel          string
	// // VideoBitRate     string
	// VideoQuality int
	// // 	Preset           string
	// // 	Scale            string
	// // 	KeyframeInterval int
}

// GetID _
func (cdr *ConvertDealerRequest) GetID() string {
	return cdr.Identity
}

// GetType _
func (cdr *ConvertDealerRequest) GetType() string {
	return ConvertV1Type
}
