package models

// SegmentMutation _
type SegmentMutation struct {
}

// NewSegmentMutation _
func NewSegmentMutation() *SegmentMutation {
	return &SegmentMutation{}
}

// CancelSegment _
func (si *SegmentMutation) CancelSegment(segment ISegment, reason string) error {
	segment.cancel(reason)
	return nil
}

// FailSegment _
func (si *SegmentMutation) FailSegment(segment ISegment, err error) error {
	segment.setLastError(err)
	segment.incrementRetriesCount()

	if segment.GetRetriesCount() < MaxRetriesCount {
		return nil
	}

	return si.CancelSegment(segment, CancellationReasonFailed)
}

// PublishSegment _
func (si *SegmentMutation) PublishSegment(segment ISegment) error {
	segment.publish()
	return nil
}

// FinishSegment _
func (si *SegmentMutation) FinishSegment(segment ISegment) error {
	segment.finish()
	return nil
}

// // RepublishSegment _
// func (si *SegmentInteraction) RepublishSegment(segment ISegment) error {
// 	segment.publish()
// 	return nil
// }

// LockSegment _
func (si *SegmentMutation) LockSegment(segment ISegment, performer IAuthor) error {
	segment.lock(performer)

	return nil
}

// UnlockSegment _
func (si *SegmentMutation) UnlockSegment(segment ISegment) error {
	segment.unlock()
	return nil
}
