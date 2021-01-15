package worker

import (
	"github.com/machinebox/progress"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/convert"
)

// ProgressMessage _
type ProgressMessage struct {
	step    models.ProgressStep
	percent float64
}

func makeProgresserFromConvert(bpm convert.BatchProgressMessage) models.Progresser {
	return &ProgressMessage{
		step:    models.ProcessingStep,
		percent: bpm.Progress.Percent(),
	}
}

func makeIoProgresser(iop progress.Progress, step models.ProgressStep) models.Progresser {
	return &ProgressMessage{
		step:    step,
		percent: iop.Percent(),
	}
}

// Step _
func (wp *ProgressMessage) Step() models.ProgressStep {
	return wp.step
}

// Percent _
func (wp *ProgressMessage) Percent() float64 {
	return wp.percent
}
