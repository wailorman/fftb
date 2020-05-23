package ffchunker

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/files"
)

// DurationCalculator _
type DurationCalculator struct {
}

// VideoDurationCalculator _
type VideoDurationCalculator interface {
	Calculate(file files.Filer) (float64, error)
}

// NewDurationCalculator _
func NewDurationCalculator() *DurationCalculator {
	return &DurationCalculator{}
}

// Calculate _
func (d *DurationCalculator) Calculate(file files.Filer) (float64, error) {
	cmdStr := fmt.Sprintf("ffprobe -i \"%s\" -show_entries format=duration -v quiet -of default=noprint_wrappers=1:nokey=1", file.FullPath())

	out, err := exec.Command("bash", "-c", cmdStr).Output()

	if err != nil {
		return 0, errors.Wrap(err, "ffprobe executing problem")
	}

	floatValue, err := strconv.ParseFloat(
		strings.ReplaceAll(string(out), "\n", ""),
		64,
	)

	if err != nil {
		return 0, errors.Wrap(err, "Converting duration to float")
	}

	return floatValue, nil
}
