package models

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test__WrapOzzoValidationError__ozzoError(t *testing.T) {
	ozzoError := validation.Errors{
		"id": errors.New("blank"),
	}

	wrappedOzzo := WrapOzzoValidationError(ozzoError)

	assert.Equal(t, true, errors.Is(wrappedOzzo, ErrInvalid), "errors.Is expectation for ozzo error")

	var iErrInvalid ValidationError
	isValidationErr := errors.As(wrappedOzzo, &iErrInvalid)

	assert.Equal(t, true, isValidationErr, "Type expectation for ozzo error")

	if isValidationErr {
		assert.Equal(t,
			map[string]string{"id": "blank"},
			iErrInvalid.Errors())
	}
}

func Test__WrapOzzoValidationError__regularError(t *testing.T) {
	regularError := errors.New("some err")
	wrappedRegularError := WrapOzzoValidationError(regularError)

	assert.Equal(t,
		true,
		errors.Is(wrappedRegularError,
			ErrInvalid),
		"errors.Is expectation for regular error")

	var iErrInvalid ValidationError
	isValidationErr := errors.As(wrappedRegularError, &iErrInvalid)

	assert.Equal(t, false, isValidationErr, "Type expectation for regular error")
}

func Test__WrapOzzoValidationError__nil(t *testing.T) {
	wrappedRegularError := WrapOzzoValidationError(nil)

	assert.Equal(t,
		true,
		errors.Is(wrappedRegularError,
			ErrInvalid),
		"errors.Is expectation for nil error")

	var iErrInvalid ValidationError
	isValidationErr := errors.As(wrappedRegularError, &iErrInvalid)

	assert.Equal(t, false, isValidationErr, "Type expectation for nil error")
}
