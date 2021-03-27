package handlers

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

// Action4TimeshiftRegexpFull _
var Action4TimeshiftRegexpFull = regexp.MustCompile("(?i)(Timeshift )(\\d{2}-\\d{2}-\\d{4} \\d{2}-\\d{2}-\\d{2})")

// Action4TimeshiftRegexp _
var Action4TimeshiftRegexp = regexp.MustCompile("(\\d{2}-\\d{2}-\\d{4} \\d{2}-\\d{2}-\\d{2})")

// Action4TimeshiftTimeLayout _
const Action4TimeshiftTimeLayout = "02-01-2006 15-04-05"

// Action4Timeshift _
type Action4Timeshift struct {
	durationCalculator DurationCalculator
}

// NewAction4Timeshift _
func NewAction4Timeshift(durationCalculator DurationCalculator) *Action4Timeshift {
	return &Action4Timeshift{
		durationCalculator: durationCalculator,
	}
}

// IsMatch _
func (gf *Action4Timeshift) IsMatch(file files.Filer) bool {
	return len(Action4TimeshiftRegexpFull.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *Action4Timeshift) Extract(file files.Filer) (time.Time, error) {
	if gf.IsMatch(file) {
		str := Action4TimeshiftRegexp.FindString(file.Name())

		videoDurationSecs, err := gf.durationCalculator.CalculateDuration(file)

		if err != nil {
			return time.Time{}, errors.Wrap(err, "Video duration calculation")
		}

		parsedTime, err := time.ParseInLocation(
			Action4TimeshiftTimeLayout,
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
func (gf *Action4Timeshift) HandlerName() string {
	return "action4_timeshift"
}
