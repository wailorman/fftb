package models

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// OrderMutator _
type OrderMutator struct {
	logger logrus.FieldLogger
}

// NewOrderMutation _
func NewOrderMutation(logger logrus.FieldLogger) *OrderMutator {
	return &OrderMutator{logger: logger}
}

// CancelOrder _
func (oi *OrderMutator) CancelOrder(segmentMutator SegmentCanceller, order IOrder, segments []ISegment, reason string) error {
	for _, segment := range segments {
		err := segmentMutator.CancelSegment(segment, CancellationReasonOrderCancelled)

		if err != nil {
			return errors.Wrapf(err, "Cancel segment `%s`", segment.GetID())
		}
	}

	order.cancel(reason)
	return nil
}

// FailOrder _
func (oi *OrderMutator) FailOrder(segmentMutator SegmentCanceller, order IOrder, segments []ISegment, err error) error {
	order.setLastError(err)
	order.incrementRetriesCount()

	if order.GetRetriesCount() < MaxRetriesCount {
		return nil
	}

	return oi.CancelOrder(segmentMutator, order, segments, CancellationReasonFailed)
}
