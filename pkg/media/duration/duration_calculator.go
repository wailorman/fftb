package cut

import (
	"github.com/pkg/errors"
	"github.com/wailorman/chunky/pkg/files"
	mediaInfo "github.com/wailorman/chunky/pkg/media/info"
)

// CalculatorInstance _
type CalculatorInstance struct {
	infoGetter mediaInfo.Getter
}

// Calculator _
type Calculator interface {
	Calculate(file files.Filer) (float64, error)
}

// NewCalculator _
func NewCalculator(infoGetter mediaInfo.Getter) *CalculatorInstance {
	return &CalculatorInstance{
		infoGetter: infoGetter,
	}
}

// Calculate _
func (dc *CalculatorInstance) Calculate(file files.Filer) (float64, error) {
	metadata, err := dc.infoGetter.GetMediaInfo(file)

	if err != nil {
		return 0, errors.Wrap(err, "Getting Media info error")
	}

	if len(metadata.Streams) == 0 {
		return 0, errors.New("No streams in media file")
	}

	return metadata.Streams[0].DurationFloat, nil
}
