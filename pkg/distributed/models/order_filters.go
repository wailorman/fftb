package models

// IOrderSearchCriteria _
type IOrderSearchCriteria interface {
	Select(order IOrder) bool
}

// EmptyOrderFilters _
func EmptyOrderFilters() IOrderSearchCriteria {
	return MergeOrderFilters()
}

// ChainOrderFilters _
func ChainOrderFilters(crits ...IOrderSearchCriteria) func(order IOrder) bool {
	return func(order IOrder) bool {
		allPassed := true

		for _, crit := range crits {
			if crit == nil {
				return false
			}

			if !crit.Select(order) {
				allPassed = false
				break
			}
		}

		return allPassed
	}
}

// MergeOrderFilters _
func MergeOrderFilters(crits ...IOrderSearchCriteria) IOrderSearchCriteria {
	return &MergedOrderSearchCriteria{
		crits: crits,
	}
}

// MergedOrderSearchCriteria _
type MergedOrderSearchCriteria struct {
	crits []IOrderSearchCriteria
}

// Select _
func (crit *MergedOrderSearchCriteria) Select(order IOrder) bool {
	allPassed := true

	for _, crit := range crit.crits {
		if crit == nil {
			return false
		}

		if !crit.Select(order) {
			allPassed = false
			break
		}
	}

	return allPassed
}

// OrderStateCriteria _
type OrderStateCriteria struct {
	state string
}

// OrderStateFilter _
func OrderStateFilter(state string) *OrderStateCriteria {
	return &OrderStateCriteria{
		state: state,
	}
}

// Select _
func (crit *OrderStateCriteria) Select(order IOrder) bool {
	return order.GetState() == crit.state
}
