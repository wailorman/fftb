package handlers

import "github.com/wailorman/ffchunker/pkg/files"

type durationCalculatorStub struct {
	value float64
}

func newDurationCalculatorStub(value float64) *durationCalculatorStub {
	return &durationCalculatorStub{
		value: value,
	}
}

// Calculate _
func (d *durationCalculatorStub) Calculate(file files.Filer) (float64, error) {
	return d.value, nil
}
