package handlers

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

// GeforceDVRRegexpFull _
var GeforceDVRRegexpFull = regexp.MustCompile("(\\d{4}\\.\\d{2}\\.\\d{2} - \\d{2}\\.\\d{2}\\.\\d{2}\\.\\d{2})(\\.DVR)")

// GeforceDVRRegexp _
var GeforceDVRRegexp = regexp.MustCompile("(\\d{4}\\.\\d{2}\\.\\d{2} - \\d{2}\\.\\d{2}\\.\\d{2}\\.\\d{2})")

// GeforceDVRTimeLayout _
const GeforceDVRTimeLayout = "2006.01.02 - 15.04.05.99"

// DurationCalculator interface _
type DurationCalculator interface {
	Calculate(file files.Filer) (float64, error)
}

// GeforceDVR _
type GeforceDVR struct {
	durationCalculator DurationCalculator
}

// NewGeforceDVR _
func NewGeforceDVR(durationCalculator DurationCalculator) *GeforceDVR {
	return &GeforceDVR{
		durationCalculator: durationCalculator,
	}
}

// IsMatch _
func (gf *GeforceDVR) IsMatch(file files.Filer) bool {
	return len(GeforceDVRRegexpFull.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *GeforceDVR) Extract(file files.Filer) (time.Time, error) {
	if gf.IsMatch(file) {
		str := GeforceDVRRegexp.FindString(file.Name())

		videoDurationSecs, err := gf.durationCalculator.Calculate(file)

		if err != nil {
			return time.Time{}, errors.Wrap(err, "Video duration calculation")
		}

		parsedTime, err := time.ParseInLocation(
			GeforceDVRTimeLayout,
			str,
			time.Now().Location(),
		)

		if err == nil {
			return parsedTime.Add(time.Duration(-videoDurationSecs) * time.Second), nil
		}
	}

	return time.Time{}, ErrNoTimeMatches
}

// HandlerName _
func (gf *GeforceDVR) HandlerName() string {
	return "geforce_dvr"
}
