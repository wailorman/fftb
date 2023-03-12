package throttle

import (
	"sync"
	"time"
)

type Throttler func(f func())

// New _
func New(every time.Duration) func(f func()) {
	d := &throttler{every: every}

	return func(f func()) {
		d.add(f)
	}
}

type throttler struct {
	lock      sync.Mutex
	every     time.Duration
	callAfter *time.Time
}

func (d *throttler) add(f func()) {
	d.lock.Lock()
	defer d.lock.Unlock()

	now := time.Now()

	if d.callAfter == nil || (d.callAfter != nil && now.After(*d.callAfter)) {
		callAfter := now.Add(d.every)
		d.callAfter = &callAfter
		f()
	}

	// if d.lastCalled != nil && d.callAfter != nil {
	// 	callAfter = *d.callAfter

	// 	if now.After(callAfter) {
	// 		callAfter := now.Add(d.every)
	// 		d.lastCalled = &now
	// 		d.callAfter = &callAfter
	// 		f()
	// 	}
	// } else {
	// 	callAfter := now.Add(d.every)
	// 	d.lastCalled = &now
	// 	d.callAfter = &callAfter
	// 	f()
	// }
}
