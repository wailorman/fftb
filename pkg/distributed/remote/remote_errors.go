package remote

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	dealerSchema "github.com/wailorman/fftb/pkg/distributed/remote/schema/dealer"
)

func parseError(clientErr error, httpResponse *http.Response, rawBody []byte, details ...*dealerSchema.ProblemDetails) error {
	if clientErr != nil {
		return errors.Wrap(clientErr, "API request failed")
	}

	if httpResponse == nil {
		return models.ErrUnknown
	}

	var targetDetails *dealerSchema.ProblemDetails

	for _, detail := range details {
		if detail != nil {
			targetDetails = detail
			break
		}
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		if targetDetails != nil {
			return errors.Wrapf(models.ErrNotFound, errCtx(httpResponse, targetDetails))
		}

		return models.ErrNotFound
	}

	if targetDetails == nil {
		if httpResponse.StatusCode >= 400 {
			return errors.Wrapf(models.ErrUnknown, "Received body: `%s`", rawBody)
		}

		return nil
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		var cause error

		switch targetDetails.Title {
		default:
			cause = getKnownError(targetDetails.Title)
		}

		return errors.Wrapf(cause, errCtx(httpResponse, targetDetails))
	}

	if httpResponse.StatusCode == http.StatusUnprocessableEntity {
		var cause error

		switch targetDetails.Title {
		case models.ErrInvalid.Error():
			if targetDetails.Fields != nil && *&targetDetails.Fields.AdditionalProperties != nil {
				cause = models.NewValidationError(*&targetDetails.Fields.AdditionalProperties)
			} else {
				cause = models.ErrInvalid
			}

		default:
			cause = getKnownError(targetDetails.Title)
		}

		return errors.Wrapf(cause, errCtx(httpResponse, targetDetails))
	}

	return errors.Wrapf(models.ErrUnknown, *targetDetails.Detail)
}

func errCtx(httpResponse *http.Response, details *dealerSchema.ProblemDetails) string {
	infoFormat := "HTTP %d - `%s` (`%s`, description: `%s`)"

	detail := ""
	eType := ""

	if details.Detail != nil {
		detail = string(*details.Detail)
	}

	if details.Type != nil {
		eType = string(*details.Type)
	}

	errCtxValues := []string{
		details.Title,
		detail,
		eType}

	return fmt.Sprintf(infoFormat, httpResponse.StatusCode, errCtxValues)
}

func getKnownError(errStr string) error {
	switch errStr {
	case models.ErrUnknownType.Error():
		return models.ErrUnknownType

	case models.ErrUnknownStorageClaimType.Error():
		return models.ErrUnknownStorageClaimType

	case models.ErrMissingStorageClaim.Error():
		return models.ErrMissingStorageClaim

	case models.ErrMissingRequest.Error():
		return models.ErrMissingRequest

	case models.ErrTimeoutReached.Error():
		return models.ErrTimeoutReached

	case models.ErrLockTimeout.Error():
		return models.ErrLockTimeout

	case models.ErrLockTimeoutReached.Error():
		return models.ErrLockTimeoutReached

	case models.ErrMissingLockAuthor.Error():
		return models.ErrMissingLockAuthor

	case models.ErrSegmentIsLocked.Error():
		return models.ErrSegmentIsLocked

	case models.ErrMissingSegment.Error():
		return models.ErrMissingSegment

	case models.ErrMissingOrder.Error():
		return models.ErrMissingOrder

	case models.ErrMissingPublisher.Error():
		return models.ErrMissingPublisher

	case models.ErrMissingPerformer.Error():
		return models.ErrMissingPerformer

	case models.ErrPerformerMismatch.Error():
		return models.ErrPerformerMismatch

	case models.ErrMissingAuthor.Error():
		return models.ErrMissingAuthor

	case models.ErrUnauthorized.Error():
		return models.ErrUnauthorized

	case models.ErrInvalidAuthorityKey.Error():
		return models.ErrInvalidAuthorityKey

	case models.ErrInvalidSessionKey.Error():
		return models.ErrInvalidSessionKey

	case models.ErrMissingAccessToken.Error():
		return models.ErrMissingAccessToken

	case models.ErrNotFound.Error():
		return models.ErrNotFound

	case models.ErrStorageClaimAlreadyAllocated.Error():
		return models.ErrStorageClaimAlreadyAllocated

	default:
		return models.ErrUnknown
	}
}
