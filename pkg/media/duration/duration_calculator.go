package cut

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/minfo"
)

// CalculatorInstance _
type CalculatorInstance struct {
	infoGetter minfo.Getter
}

// Calculator _
type Calculator interface {
	CalculateDuration(files.Filer) (float64, error)
}

// NewCalculator _
func NewCalculator(infoGetter minfo.Getter) *CalculatorInstance {
	return &CalculatorInstance{
		infoGetter: infoGetter,
	}
}

// CalculateDuration _
func (dc *CalculatorInstance) CalculateDuration(file files.Filer) (float64, error) {
	metadata, err := dc.infoGetter.GetMediaInfo(file)

	if err != nil {
		return 0, errors.Wrap(err, "Getting Media info error")
	}

	if len(metadata.Streams) == 0 {
		return 0, errors.New("No streams in media file")
	}

	return metadata.Streams[0].DurationFloat, nil
}
