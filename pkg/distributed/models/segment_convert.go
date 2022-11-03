package models

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// SegmentStatePrepared _
const SegmentStatePrepared = "prepared"

// SegmentStatePublished _
const SegmentStatePublished = "published"

// SegmentStateInProgress is dynamic state, used only in some presenters.
// Can be returned by GetCurrentState()
const SegmentStateInProgress = "in_progress"

// SegmentStateAccepted _
const SegmentStateAccepted = "accepted"

// SegmentStateFinished _
const SegmentStateFinished = "finished"

// SegmentStateCancelled _
const SegmentStateCancelled = "cancelled"

// ConvertSegment _
type ConvertSegment struct {
	Identity                   string `json:"identity"`
	OrderIdentity              string `json:"order_identity"` // TODO: Deprecated
	Type                       string `json:"type"`
	InputStorageClaimIdentity  string `json:"input_storage_claim_identity"`  // TODO: Deprecated
	OutputStorageClaimIdentity string `json:"output_storage_claim_identity"` // TODO: Deprecated
	State                      string `json:"state"`                         // TODO: Deprecated

	Params   convert.Params `json:"params"`
	Muxer    string         `json:"muxer"`
	Position int            `json:"position"` // TODO: Deprecated

	Publisher   IAuthor    `json:"publisher"`    // TODO: Deprecated
	LockedUntil *time.Time `json:"locked_until"` // TODO: Deprecated
	LockedBy    IAuthor    `json:"locked_by"`    // TODO: Deprecated

	RetriesCount       int        `json:"retries_count"`       // TODO: Deprecated
	RetryAt            *time.Time `json:"retry_at"`            // TODO: Deprecated
	LastError          string     `json:"last_error"`          // TODO: Deprecated
	CancellationReason string     `json:"cancellation_reason"` // TODO: Deprecated
}

// GetID _
func (ct *ConvertSegment) GetID() string {
	return ct.Identity
}

// GetType _
func (ct *ConvertSegment) GetType() string {
	return ConvertV1Type
}

// GetOrderID _
func (ct *ConvertSegment) GetOrderID() string {
	return ct.OrderIdentity
}

// GetInputStorageClaimIdentity _
func (ct *ConvertSegment) GetInputStorageClaimIdentity() string {
	return ct.InputStorageClaimIdentity
}

// GetOutputStorageClaimIdentity _
func (ct *ConvertSegment) GetOutputStorageClaimIdentity() string {
	return ct.OutputStorageClaimIdentity
}

// GetPayload _
func (ct *ConvertSegment) GetPayload() (string, error) {
	b, err := json.Marshal(ct)

	return string(b), err
}

// GetIsLocked _
func (ct *ConvertSegment) GetIsLocked() bool {
	if ct.LockedUntil == nil || ct.LockedBy == nil {
		return false
	}

	return !time.Now().After(*ct.LockedUntil)
}

// GetLockedBy _
func (ct *ConvertSegment) GetLockedBy() IAuthor {
	if !ct.GetIsLocked() {
		return nil
	}

	return ct.LockedBy
}

// GetLockedUntil _
func (ct *ConvertSegment) GetLockedUntil() *time.Time {
	if !ct.GetIsLocked() {
		return nil
	}

	return ct.LockedUntil
}

// // GetStorageClaim _
// func (ct *ConvertSegment) GetStorageClaim() IStorageClaim {
// 	return ct.StorageClaim
// }

// Failed _
// func (ct *ConvertSegment) Failed(err error) {
// 	// TODO:
// 	// panic(ErrNotImplemented)
// 	panic(err)
// 	// return
// }

// GetState _
func (ct *ConvertSegment) GetState() string {
	return ct.State
}

func (ct *ConvertSegment) cancel(reason string) {
	ct.State = SegmentStateCancelled
	ct.CancellationReason = reason
	ct.unlock()
}

func (ct *ConvertSegment) publish() {
	ct.State = SegmentStatePublished
	ct.unlock()
}

func (ct *ConvertSegment) finish() {
	ct.State = SegmentStateFinished
	ct.unlock()
}

// GetCurrentState _
func (ct *ConvertSegment) GetCurrentState() string {
	if ct.GetIsLocked() {
		return SegmentStateInProgress
	}

	return ct.State
}

// GetPublisher _
func (ct *ConvertSegment) GetPublisher() IAuthor {
	return ct.Publisher
}

// GetPerformer _
func (ct *ConvertSegment) GetPerformer() IAuthor {
	return ct.LockedBy
}

// MatchPublisher _
func (ct *ConvertSegment) MatchPublisher(publisher IAuthor) bool {
	if ct.Publisher == nil {
		return false
	}

	if publisher == nil {
		return false
	}

	return ct.Publisher == publisher
}

// MatchPerformer _
func (ct *ConvertSegment) MatchPerformer(performer IAuthor) bool {
	if ct.LockedBy == nil {
		return false
	}

	if !ct.GetIsLocked() {
		return false
	}

	if performer == nil {
		return false
	}

	return ct.LockedBy == performer
}

func (ct *ConvertSegment) lock(performer IAuthor) {
	lockedUntil := time.Now().Add(SegmentLockDuration)

	ct.LockedBy = performer
	ct.LockedUntil = &lockedUntil
}

func (ct *ConvertSegment) unlock() {
	ct.LockedBy = nil
	ct.LockedUntil = nil
}

// GetPosition _
func (ct *ConvertSegment) GetPosition() int {
	return ct.Position
}

// GetRetriesCount _
func (ct *ConvertSegment) GetRetriesCount() int {
	return ct.RetriesCount
}

// GetRetryAt _
func (ct *ConvertSegment) GetRetryAt() *time.Time {
	return ct.RetryAt
}

// GetCanRetry _
func (ct *ConvertSegment) GetCanRetry() bool {
	if ct.GetRetriesCount() >= MaxRetriesCount {
		return false
	}

	if ct.GetRetryAt() != nil {
		return time.Now().After(*ct.GetRetryAt())
	}

	return true
}

func (ct *ConvertSegment) incrementRetriesCount() {
	ct.RetriesCount++
	nextRetry := time.Now().Add(NextRetryOffset)
	ct.RetryAt = &nextRetry
}

func (ct *ConvertSegment) setLastError(err error) {
	ct.LastError = err.Error()
}

// GetCanPerform _
func (ct *ConvertSegment) GetCanPerform() bool {
	canRetry := true

	if ct.GetRetryAt() != nil {
		canRetry = time.Now().After(*ct.GetRetryAt())
	}

	return !ct.GetIsLocked() &&
		ct.GetRetriesCount() < MaxRetriesCount &&
		canRetry &&
		ct.GetState() == SegmentStatePublished

}

// Validate _
func (ct ConvertSegment) Validate() error {
	validators := make([]*validation.FieldRules, 0)

	validators = append(validators,
		validation.Field(&ct.State,
			validation.Required,
			validation.In(
				SegmentStateCancelled,
				SegmentStatePrepared,
				SegmentStatePublished,
				SegmentStateAccepted,
				SegmentStateFinished)))

	if ct.State == SegmentStatePublished {
		validators = append(validators,
			validation.Field(&ct.InputStorageClaimIdentity, validation.Required))
	}

	if ct.State == SegmentStateFinished {
		validators = append(validators,
			validation.Field(&ct.OutputStorageClaimIdentity, validation.Required))
	}

	if ct.State == SegmentStateCancelled {
		validators = append(validators,
			validation.Field(&ct.CancellationReason,
				validation.Required))
	}

	validators = append(validators,
		validation.Field(&ct.Type, validation.Required, validation.In(ConvertV1Type)),
		validation.Field(&ct.Identity, validation.Required),
		validation.Field(&ct.OrderIdentity, validation.Required),
		validation.Field(&ct.Muxer, validation.Required),
		validation.Field(&ct.Publisher, validation.Required))

	return WrapOzzoValidationError(validation.ValidateStruct(&ct, validators...))
}
