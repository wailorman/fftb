package dlog

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// KeyOrderID _
const KeyOrderID = "order_id"

// KeySegmentID _
const KeySegmentID = "segment_id"

// KeySegmentState _
const KeySegmentState = "segment_state"

// KeyPercent _
const KeyPercent = "percent"

// KeyPerformer _
const KeyPerformer = "performer"

// KeyStorePayload _
const KeyStorePayload = "store_payload"

// KeyCallee _
const KeyCallee = "callee"

// KeyStorageClaim _
const KeyStorageClaim = "storage_claim"

// KeyReason _
const KeyReason = "reason"

// PrefixContracterPublishWorker _
const PrefixContracterPublishWorker = "fftb.contracter.publish_worker"

// PrefixContracterConcatWorker _
const PrefixContracterConcatWorker = "fftb.contracter.concat_worker"

// PrefixContracter _
const PrefixContracter = "fftb.contracter"

// PrefixDealer _
const PrefixDealer = "fftb.dealer"

// PrefixWorker _
const PrefixWorker = "fftb.worker"

// PrefixSegmConcatOperation _
const PrefixSegmConcatOperation = "fftb.segm.concat_operation"

// PrefixSegmSliceOperation _
const PrefixSegmSliceOperation = "fftb.segm.slice_operation"

// PrefixAPI _
const PrefixAPI = "fftb.api"

// SegmentProgress _
func SegmentProgress(logger logrus.FieldLogger, seg models.ISegment, p models.IProgress) {
	entry := WithSegment(logger, seg).
		WithField(KeyPercent, fmt.Sprintf("%.4f", p.Percent())).
		WithField(KeySegmentState, p.Step())

	entry.Debug("Processing segment")
}

// BasicProgress _
type BasicProgress struct {
	step    models.ProgressStep
	percent float64
}

// TODO: Remove

// BuildProgress is internal function
func BuildProgress(step models.ProgressStep, percent float64) *BasicProgress {
	return &BasicProgress{step, percent}
}

// Step _
func (bP *BasicProgress) Step() models.ProgressStep {
	return bP.step
}

// Percent _
func (bP *BasicProgress) Percent() float64 {
	return bP.percent
}

// WithOrder _
func WithOrder(logger logrus.FieldLogger, order models.IOrder) logrus.FieldLogger {
	return logger.WithField(KeyOrderID, order.GetID())
}

// WithSegment _
func WithSegment(logger logrus.FieldLogger, segment models.ISegment) logrus.FieldLogger {
	return logger.WithField(KeyOrderID, segment.GetOrderID()).
		WithField(KeySegmentID, segment.GetID())
}

// MakeIOProgress _
func MakeIOProgress(step models.ProgressStep, percent float64) models.IProgress {
	return &BasicProgress{
		step:    step,
		percent: percent,
	}
}

// JSON _
func JSON(v interface{}) string {
	if v == nil {
		return "{}"
	}

	bytes, err := json.Marshal(v)

	if err != nil {
		return fmt.Sprintf("<failed to generate json: %s>", err)
	}

	return string(bytes)
}
