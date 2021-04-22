package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/schema"
)

// APIError _
type APIError struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func newAPIError(err error) (code int, body *schema.ProblemDetails) {
	cause := errors.Cause(err)

	code = http.StatusUnprocessableEntity

	if errors.Is(err, models.ErrNotFound) {
		code = http.StatusNotFound
	}

	if errors.Is(err, models.ErrUnknownType) {
		code = http.StatusUnprocessableEntity
	}

	if errors.Is(err, models.ErrUnknown) {
		code = http.StatusInternalServerError
	}

	if errors.Is(err, models.ErrMissingAccessToken) ||
		errors.Is(err, models.ErrInvalidSessionKey) ||
		errors.Is(err, models.ErrInvalidAuthorityKey) ||
		errors.Is(err, models.ErrMissingAuthor) {

		code = http.StatusUnauthorized
	}

	detail := err.Error()
	errType := "github.com/wailorman/fftb"

	problemDetails := &schema.ProblemDetails{
		Type:   &errType,
		Title:  cause.Error(),
		Detail: &detail,
	}

	var validationErr models.ValidationError
	if errors.As(err, &validationErr) {
		problemDetails.Title = models.ErrInvalid.Error()
		problemDetails.Fields = &schema.ProblemDetails_Fields{}

		for key, val := range validationErr.Errors() {
			problemDetails.Fields.Set(key, val)
		}
	}

	var echoError *echo.HTTPError
	if errors.As(err, &echoError) {
		code = echoError.Code
		problemDetails.Title = models.ErrUnknown.Error()
	}

	return code, problemDetails
}
