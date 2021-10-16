package models

// PercentProgress implements models.Progresser interface
type PercentProgress struct {
	percent float64
}

// NewPercentProgress creates new simple PercentProgress object
func NewPercentProgress(percent float64) *PercentProgress {
	return &PercentProgress{percent}
}

// Percent _
func (pp *PercentProgress) Percent() float64 {
	return pp.percent
}

// Step _
func (pp *PercentProgress) Step() ProgressStep {
	return ProcessingStep
}
