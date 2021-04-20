package models

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ConvertOrder _
type ConvertOrder struct {
	Identity  string         `json:"identity"`
	Type      string         `json:"type"`
	State     string         `json:"state"`
	InFile    files.Filer    `json:"in_file"`
	OutFile   files.Filer    `json:"out_file"`
	Params    convert.Params `json:"params"`
	Publisher IAuthor        `json:"publisher"`

	RetriesCount       int        `json:"retries_count"`
	RetryAt            *time.Time `json:"retry_at"`
	LastError          string     `json:"last_error"`
	CancellationReason string     `json:"cancellation_reason"`
}

// OrderStateQueued _
const OrderStateQueued = "queued"

// OrderStateInProgress _
const OrderStateInProgress = "in_progress"

// OrderStateFinished _
const OrderStateFinished = "finished"

// OrderStateCancelled _
const OrderStateCancelled = "cancelled"

// GetID _
func (co *ConvertOrder) GetID() string {
	return co.Identity
}

// // SetID _
// func (co *ConvertOrder) SetID(id string) {
// 	co.Identity = id
// }

// GetType _
func (co *ConvertOrder) GetType() string {
	return ConvertV1Type
}

// // SetType _
// func (co *ConvertOrder) SetType(orderType string) {
// 	co.Type = orderType
// }

// GetPayload _
func (co *ConvertOrder) GetPayload() (string, error) {
	b, err := json.Marshal(co)

	return string(b), err
}

// GetState _
func (co *ConvertOrder) GetState() string {
	return co.State
}

func (co *ConvertOrder) cancel(reason string) {
	co.State = OrderStateCancelled
	co.CancellationReason = reason
}

// GetPublisher _
func (co *ConvertOrder) GetPublisher() IAuthor {
	return co.Publisher
}

// MatchPublisher _
func (co *ConvertOrder) MatchPublisher(publisher IAuthor) bool {
	if co.Publisher == nil {
		return false
	}

	if publisher == nil {
		return false
	}

	return co.Publisher.IsEqual(publisher)
}

// CalculateProgress _
func (co *ConvertOrder) CalculateProgress(segments []ISegment) float64 {
	if len(segments) == 0 {
		return 0
	}

	totalSegments := float64(len(segments))
	finishedSegments := 0.0

	for _, segment := range segments {
		if segment.GetState() == SegmentStateFinished {
			finishedSegments++
		}
	}

	return finishedSegments / totalSegments
}

// GetInputFile _
func (co *ConvertOrder) GetInputFile() files.Filer {
	return co.InFile
}

// GetOutputFile _
func (co *ConvertOrder) GetOutputFile() files.Filer {
	return co.OutFile
}

// GetRetriesCount _
func (co *ConvertOrder) GetRetriesCount() int {
	return co.RetriesCount
}

// GetRetryAt _
func (co *ConvertOrder) GetRetryAt() *time.Time {
	return co.RetryAt
}

// GetCanRetry _
func (co *ConvertOrder) GetCanRetry() bool {
	if co.GetRetriesCount() >= MaxRetriesCount {
		return false
	}

	if co.GetRetryAt() != nil {
		return time.Now().After(*co.GetRetryAt())
	}

	return true
}

// GetCanPublish _
func (co *ConvertOrder) GetCanPublish() bool {
	return co.GetState() == OrderStateQueued && co.GetCanRetry()
}

// GetCanConcat _
func (co *ConvertOrder) GetCanConcat(segments []ISegment) bool {
	allSegmentsFinished := true

	for _, segment := range segments {
		if segment.GetState() != SegmentStateFinished {
			allSegmentsFinished = false
			break
		}
	}

	return allSegmentsFinished && co.GetState() == OrderStateInProgress && co.GetCanRetry()
}

func (co *ConvertOrder) incrementRetriesCount() {
	co.RetriesCount++
	nextRetry := time.Now().Add(NextRetryOffset)
	co.RetryAt = &nextRetry
}

func (co *ConvertOrder) setLastError(err error) {
	co.LastError = err.Error()
}

// Validate validates convert order object and returns ValidationError or nil
func (co ConvertOrder) Validate() error {
	validators := make([]*validation.FieldRules, 0)

	validators = append(validators,
		validation.Field(&co.State,
			validation.Required,
			validation.In(
				OrderStateCancelled,
				OrderStateQueued,
				OrderStateInProgress,
				OrderStateFinished)))

	if co.State == OrderStateCancelled {
		validators = append(validators,
			validation.Field(&co.CancellationReason,
				validation.Required))
	}

	validators = append(validators,
		validation.Field(&co.Type, validation.Required, validation.In(ConvertV1Type)),
		validation.Field(&co.Identity, validation.Required),
		validation.Field(&co.InFile, validation.Required),
		validation.Field(&co.OutFile, validation.Required),
		validation.Field(&co.Publisher, validation.Required))

	return WrapOzzoValidationError(validation.ValidateStruct(&co, validators...))
}
