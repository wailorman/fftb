package chwg

import "sync"

// New _
func New() *ChannelledWaitGroup {
	return &ChannelledWaitGroup{
		wg: &sync.WaitGroup{},
	}
}

// ChannelledWaitGroup _
type ChannelledWaitGroup struct {
	counter int
	wg      *sync.WaitGroup
}

// Add _
func (cwg *ChannelledWaitGroup) Add(delta int) {
	cwg.counter += delta
	cwg.wg.Add(delta)
}

// Done _
func (cwg *ChannelledWaitGroup) Done() {
	cwg.counter--
	cwg.wg.Done()
}

// AllDone _
func (cwg *ChannelledWaitGroup) AllDone() {
	for i := 0; i <= cwg.counter; i++ {
		cwg.Done()
	}
}

// IsFinished _
func (cwg *ChannelledWaitGroup) IsFinished() bool {
	return cwg.counter < 1
}

// Wait _
func (cwg *ChannelledWaitGroup) Wait() {
	cwg.wg.Wait()
}

// Closed _
func (cwg *ChannelledWaitGroup) Closed() <-chan struct{} {
	ch := make(chan struct{})
	cwg.wg.Wait()
	close(ch)
	return ch
}
