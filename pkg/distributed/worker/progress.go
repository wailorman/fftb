package worker

import (
	"github.com/machinebox/progress"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/ff"
)

// ProgressMessage _
type ProgressMessage struct {
	step    models.ProgressStep
	percent float64
}

func makeProgresserFromConvert(pm ff.Progressable) models.IProgress {
	return &ProgressMessage{
		step:    models.ProcessingStep,
		percent: pm.Percent(),
	}
}

func makeIoProgresser(iop progress.Progress, step models.ProgressStep) models.IProgress {
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
