package models

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

// ValidationErrorObject _
type ValidationErrorObject struct {
	errors map[string]string
}

// ValidationError _
type ValidationError interface {
	error
	Errors() map[string]string
}

// WrapOzzoValidationError _
func WrapOzzoValidationError(err error) ValidationError {
	if validationErr, isValidationErr := errors.Cause(err).(validation.Errors); isValidationErr {
		verr := &ValidationErrorObject{
			errors: make(map[string]string),
		}

		for key, vErr := range validationErr {
			verr.errors[key] = vErr.Error()
		}

		return verr
	}

	return &ValidationErrorObject{}
}

// NewValidationError _
func NewValidationError(errosMap map[string]string) ValidationError {
	if errosMap == nil {
		return &ValidationErrorObject{}
	}

	return &ValidationErrorObject{
		errors: errosMap,
	}
}

// String _
func (verr *ValidationErrorObject) String() string {
	fields := make([]string, 0)

	for key, val := range verr.errors {
		fields = append(fields, fmt.Sprintf("%s: %s", key, val))
	}

	return strings.Join(fields, "; ")
}

// Error _
func (verr *ValidationErrorObject) Error() string {
	return verr.String()
}

// Unwrap _
func (verr *ValidationErrorObject) Unwrap() error {
	return ErrInvalid
}

// Errors _
func (verr *ValidationErrorObject) Errors() map[string]string {
	return verr.errors
}
