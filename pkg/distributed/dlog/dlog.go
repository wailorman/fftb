package dlog

import (
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

// SegmentProgress _
func SegmentProgress(logger logrus.FieldLogger, seg models.ISegment, p models.Progresser) {
	entry := logger.
		WithField(KeyPercent, fmt.Sprintf("%.4f", p.Percent())).
		WithField(KeySegmentID, seg.GetID()).
		WithField(KeySegmentState, p.Step()).
		WithField(KeyOrderID, seg.GetOrderID())

	entry.Info("Processing segment")
}
