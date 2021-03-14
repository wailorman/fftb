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

// SegmentProgress _
func SegmentProgress(logger logrus.FieldLogger, seg models.ISegment, p models.Progresser) {
	entry := logger.
		WithField(KeyPercent, fmt.Sprintf("%.4f", p.Percent())).
		WithField(KeySegmentID, seg.GetID()).
		WithField(KeySegmentState, p.Step()).
		WithField(KeyOrderID, seg.GetOrderID())

	entry.Debug("Processing segment")
}

// basicProgress _
type basicProgress struct {
	step    models.ProgressStep
	percent float64
}

// Step _
func (bP *basicProgress) Step() models.ProgressStep {
	return bP.step
}

// Percent _
func (bP *basicProgress) Percent() float64 {
	return bP.percent
}

// MakeIOProgress _
func MakeIOProgress(step models.ProgressStep, percent float64) models.Progresser {
	return &basicProgress{
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
