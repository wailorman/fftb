package handlers

import (
	"regexp"
	"time"

	"github.com/wailorman/ffchunker/files"
)

// GeforceFullRegexp _
var GeforceFullRegexp = regexp.MustCompile("(\\d{4}\\.\\d{2}\\.\\d{2} - \\d{2}\\.\\d{2}\\.\\d{2}\\.\\d{2})")

// GeforceFullTimeLayout _
const GeforceFullTimeLayout = "2006.01.02 - 15.04.05.99"

// GeforceFull _
type GeforceFull struct {
}

// NewGeforceFull _
func NewGeforceFull() *GeforceFull {
	return &GeforceFull{}
}

// IsMatch _
func (gf *GeforceFull) IsMatch(file files.Filer) bool {
	return len(GeforceFullRegexp.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *GeforceFull) Extract(file files.Filer) (time.Time, error) {
	strs := GeforceFullRegexp.FindAllString(file.Name(), -1)

	if len(strs) > 0 {
		parsedTime, err := time.ParseInLocation(
			GeforceFullTimeLayout,
			strs[0],
			time.Now().Location(),
		)

		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, ErrNoTimeMatches
}

// HandlerName _
func (gf *GeforceFull) HandlerName() string {
	return "geforce_full"
}
