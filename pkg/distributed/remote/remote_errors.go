package remote

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

func parseError(clientErr error, httpResponse *http.Response, details ...*ProblemDetails) error {
	if clientErr != nil {
		return errors.Wrap(clientErr, "API request failed")
	}

	if httpResponse == nil {
		return models.ErrUnknown
	}

	var targetDetails *ProblemDetails

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
		return nil
	}

	if httpResponse.StatusCode == http.StatusUnprocessableEntity {
		var cause error

		switch targetDetails.Title {
		case models.ErrUnknownType.Error():
			cause = models.ErrUnknownType

		case models.ErrUnknownStorageClaimType.Error():
			cause = models.ErrUnknownStorageClaimType

		case models.ErrMissingStorageClaim.Error():
			cause = models.ErrMissingStorageClaim

		case models.ErrMissingRequest.Error():
			cause = models.ErrMissingRequest

		case models.ErrTimeoutReached.Error():
			cause = models.ErrTimeoutReached

		case models.ErrLockTimeout.Error():
			cause = models.ErrLockTimeout

		case models.ErrLockTimeoutReached.Error():
			cause = models.ErrLockTimeoutReached

		case models.ErrMissingLockAuthor.Error():
			cause = models.ErrMissingLockAuthor

		case models.ErrSegmentIsLocked.Error():
			cause = models.ErrSegmentIsLocked

		case models.ErrMissingSegment.Error():
			cause = models.ErrMissingSegment

		case models.ErrMissingOrder.Error():
			cause = models.ErrMissingOrder

		case models.ErrMissingPublisher.Error():
			cause = models.ErrMissingPublisher

		case models.ErrMissingPerformer.Error():
			cause = models.ErrMissingPerformer

		case models.ErrPerformerMismatch.Error():
			cause = models.ErrPerformerMismatch

		case models.ErrInvalid.Error():
			if targetDetails.Fields != nil && *&targetDetails.Fields.AdditionalProperties != nil {
				cause = models.NewValidationError(*&targetDetails.Fields.AdditionalProperties)
			} else {
				cause = models.ErrInvalid
			}

		default:
			cause = models.ErrUnknown
		}

		return errors.Wrapf(cause, errCtx(httpResponse, targetDetails))
	}

	return errors.Wrapf(models.ErrUnknown, *targetDetails.Detail)
}

func errCtx(httpResponse *http.Response, details *ProblemDetails) string {
	infoFormat := "HTTP %d - `%s` (`%s`, description: `%s`)"

	var detail, eType string

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
