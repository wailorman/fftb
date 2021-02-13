package models

import (
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

// ConvertDealerRequest _
type ConvertDealerRequest struct {
	Type          string
	Identity      string
	OrderIdentity string

	Params   convert.Params
	Muxer    string
	Position int

	Author IAuthor
}

// GetID _
func (cdr *ConvertDealerRequest) GetID() string {
	return cdr.Identity
}

// GetType _
func (cdr *ConvertDealerRequest) GetType() string {
	return ConvertV1Type
}

// GetAuthor _
func (cdr *ConvertDealerRequest) GetAuthor() IAuthor {
	return cdr.Author
}
