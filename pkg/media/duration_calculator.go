package media

import (
	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/pkg/files"
)

// DurationCalculatorInstance _
type DurationCalculatorInstance struct {
	mediaInfoGetter InfoGetter
}

// DurationCalculator _
type DurationCalculator interface {
	Calculate(file files.Filer) (float64, error)
}

// NewDurationCalculator _
func NewDurationCalculator(mediaInfoGetter InfoGetter) *DurationCalculatorInstance {
	return &DurationCalculatorInstance{
		mediaInfoGetter: mediaInfoGetter,
	}
}

// Calculate _
func (dc *DurationCalculatorInstance) Calculate(file files.Filer) (float64, error) {
	metadata, err := dc.mediaInfoGetter.GetMediaInfo(file)

	if err != nil {
		return 0, errors.Wrap(err, "Getting Media info error")
	}

	if len(metadata.Streams) == 0 {
		return 0, errors.New("No streams in media file")
	}

	return metadata.Streams[0].DurationFloat, nil
}
