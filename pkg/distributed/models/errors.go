package models

import (
	"github.com/pkg/errors"
)

// ErrUnknownType _
var ErrUnknownType = errors.New("Unknown type")

// ErrNotImplemented _
var ErrNotImplemented = errors.New("Not implemented")

// ErrUnknownStorageClaimType _
var ErrUnknownStorageClaimType = errors.New("Unknown storage claim type")

// ErrMissingStorageClaim _
var ErrMissingStorageClaim = errors.New("Missing storage claim")

// ErrMissingRequest _
var ErrMissingRequest = errors.New("Missing request")

// ErrUnknown _
var ErrUnknown = errors.New("Unknown error")

// ErrInvalid _
var ErrInvalid = errors.New("Validation error")

// ErrNotFound _
var ErrNotFound = errors.New("Not found")

// ErrTimeoutReached _
var ErrTimeoutReached = errors.New("Timeout reached")

// ErrLockTimeout _
var ErrLockTimeout = errors.New("Free segment lock timeout") // TODO: Subst. with ErrLockTimeoutReached

// ErrLockTimeoutReached _
var ErrLockTimeoutReached = errors.New("Lock timeout reached")

// ErrMissingLockAuthor _
var ErrMissingLockAuthor = errors.New("Missing lock author")

// ErrSegmentIsLocked _
var ErrSegmentIsLocked = errors.New("Segment is locked")

// ErrMissingSegment _
var ErrMissingSegment = errors.New("Missing Segment")

// ErrMissingOrder _
var ErrMissingOrder = errors.New("Missing Order")

// ErrMissingPublisher _
var ErrMissingPublisher = errors.New("Missing publisher")

// ErrMissingPerformer _
var ErrMissingPerformer = errors.New("Missing performer")

// ErrPerformerMismatch _
var ErrPerformerMismatch = errors.New("Performer mismatch")

// ErrUnauthorized _
var ErrUnauthorized = errors.New("Unauthorized")

// ErrMissingAuthor _
var ErrMissingAuthor = errors.New("Missing author") // TODO: replace above with this error

// ErrInvalidAuthorityKey _
var ErrInvalidAuthorityKey = errors.New("Invalid authority key")

// ErrInvalidSessionKey _
var ErrInvalidSessionKey = errors.New("Invalid session key")

// ErrMissingAccessToken _
var ErrMissingAccessToken = errors.New("Missing access token")

// ErrStorageClaimAlreadyAllocated _
var ErrStorageClaimAlreadyAllocated = errors.New("Storage claim already allocated")
