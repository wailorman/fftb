package remote

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
)

func makeHTTPResponse(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
	}
}

var pdType = "github.com/wailorman/fftb"
var pdDetails = "id: cannot be blank"

func Test__parseError__Validation(t *testing.T) {
	apiErr := parseError(
		noContentBody,
		nil,
		makeHTTPResponse(422),
		nil,
		nil,
		&dSchema.ProblemDetails{
			Title:  models.ErrInvalid.Error(),
			Type:   &pdType,
			Detail: &pdDetails,
			Fields: &dSchema.ProblemDetails_Fields{
				AdditionalProperties: map[string]string{
					"id": "blank",
				},
			},
		},
		nil,
	)

	var iErrInvalid models.ValidationError
	isValidationErr := errors.As(apiErr, &iErrInvalid)

	assert.Equalf(
		t,
		true,
		isValidationErr,
		"Expected validation error, received `%#v`", apiErr,
	)

	if isValidationErr {
		assert.Equal(t,
			map[string]string{"id": "blank"},
			iErrInvalid.Errors())
	}
}

func Test__parseError__UnknownType(t *testing.T) {
	apiErr := parseError(
		noContentBody,
		nil,
		makeHTTPResponse(422),
		nil,
		nil,
		&dSchema.ProblemDetails{
			Title: models.ErrUnknownType.Error(),
			Type:  &pdType,
		},
		nil,
	)

	if assert.NotNil(t, apiErr) {
		assert.Equalf(
			t,
			true,
			errors.Is(apiErr, models.ErrUnknownType),
			"Expected Unknown Type error, received `%#v`", apiErr,
		)
	}
}

func Test__parseError__NotFound(t *testing.T) {
	apiErr := parseError(noContentBody, nil, makeHTTPResponse(404), nil, nil, nil)

	if assert.NotNil(t, apiErr) {
		assert.Equalf(
			t,
			true,
			errors.Is(apiErr, models.ErrNotFound),
			"Expected not found, received `%#v`", apiErr,
		)
	}
}

func Test__parseError__Success(t *testing.T) {
	apiErr := parseError(noContentBody, nil, makeHTTPResponse(200), nil, nil, nil, nil)

	assert.Nilf(
		t,
		apiErr,
		"Expected nil, received `%#v`", apiErr,
	)
}

func Test__parseError__UnknownErr(t *testing.T) {
	apiErr := parseError(noContentBody, nil, makeHTTPResponse(500), []byte("Service unavailable"), nil, nil, nil)

	if assert.NotNil(t, apiErr) {
		assert.Equalf(
			t,
			true,
			errors.Is(apiErr, models.ErrUnknown),
			"Expected not found, received `%#v`", apiErr,
		)

		assert.Contains(
			t,
			apiErr.Error(),
			"Service unavailable",
		)
	}
}

func Test__parseError__requestErr(t *testing.T) {
	ErrSome := errors.New("SOME_ERR")

	apiErr := parseError(noContentBody, ErrSome, makeHTTPResponse(0), nil, nil, nil, nil)

	if assert.NotNil(t, apiErr) {
		assert.Equalf(
			t,
			true,
			errors.Is(apiErr, ErrSome),
			"Expected some error, received `%#v`", apiErr,
		)
	}
}
