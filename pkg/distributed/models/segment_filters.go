package models

// ISegmentSearchCriteria _
type ISegmentSearchCriteria interface {
	Select(segment ISegment) bool
}

// EmptySegmentFilters _
func EmptySegmentFilters() ISegmentSearchCriteria {
	return MergeSegmentFilters()
}

// ChainSegmentFilters _
func ChainSegmentFilters(crits ...ISegmentSearchCriteria) func(segment ISegment) bool {
	return func(segment ISegment) bool {
		allPassed := true

		for _, crit := range crits {
			if crit == nil {
				return false
			}

			if !crit.Select(segment) {
				allPassed = false
				break
			}
		}

		return allPassed
	}
}

// MergeSegmentFilters _
func MergeSegmentFilters(crits ...ISegmentSearchCriteria) ISegmentSearchCriteria {
	return &MergedSegmentSearchCriteria{
		crits: crits,
	}
}

// MergedSegmentSearchCriteria _
type MergedSegmentSearchCriteria struct {
	crits []ISegmentSearchCriteria
}

// Select _
func (crit *MergedSegmentSearchCriteria) Select(segment ISegment) bool {
	allPassed := true

	for _, crit := range crit.crits {
		if crit == nil {
			return false
		}

		if !crit.Select(segment) {
			allPassed = false
			break
		}
	}

	return allPassed
}

// SegmentStateCriteria _
type SegmentStateCriteria struct {
	state string
}

// SegmentStateFilter _
func SegmentStateFilter(state string) *SegmentStateCriteria {
	return &SegmentStateCriteria{
		state: state,
	}
}

// Select _
func (crit *SegmentStateCriteria) Select(segment ISegment) bool {
	return segment.GetState() == crit.state
}
