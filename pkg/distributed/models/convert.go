package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertV1Type _
const ConvertV1Type = "convert/v1"

// ConvertContracterRequest _
type ConvertContracterRequest struct {
	InFile  files.Filer
	OutFile files.Filer
	Params  convert.Params
	Author  IAuthor
}

// GetType _
func (cr *ConvertContracterRequest) GetType() string {
	return ConvertV1Type
}

// GetAuthor _
func (cr *ConvertContracterRequest) GetAuthor() IAuthor {
	return cr.Author
}

// Validate _
func (cr ConvertContracterRequest) Validate() error {
	verr := validation.ValidateStruct(&cr,
		validation.Field(&cr.InFile, validation.Required),
		validation.Field(&cr.OutFile, validation.Required),
		validation.Field(&cr.Author, validation.Required))

	return WrapOzzoValidationError(verr)
}

// ConvertDealerRequest _
type ConvertDealerRequest struct {
	Type          string `json:"type"`
	Identity      string `json:"id"`
	OrderIdentity string `json:"order_id"`

	Params   convert.Params `json:"params"`
	Muxer    string         `json:"muxer"`
	Position int            `json:"position"`
}

// GetID _
func (cdr *ConvertDealerRequest) GetID() string {
	return cdr.Identity
}

// GetType _
func (cdr *ConvertDealerRequest) GetType() string {
	return ConvertV1Type
}

// Validate _
func (cdr ConvertDealerRequest) Validate() error {
	verr := validation.ValidateStruct(&cdr,
		validation.Field(&cdr.Type, validation.Required, validation.In(ConvertV1Type)),
		validation.Field(&cdr.Identity, validation.Required),
		validation.Field(&cdr.OrderIdentity, validation.Required),
		validation.Field(&cdr.Muxer, validation.Required))

	return WrapOzzoValidationError(verr)
}
