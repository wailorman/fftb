package handlers

import (
	"regexp"
	"time"

	"github.com/wailorman/ffchunker/pkg/files"
)

// PlaysTvRegexp _
var PlaysTvRegexp = regexp.MustCompile("(\\d{4}_\\d{2}_\\d{2}_\\d{2}_\\d{2}_\\d{2})")

// 2016_10_19_23_23_00-ses.mp4
// 2016_10_19_23_23_00-ses

// PlaysTvTimeLayout _
const PlaysTvTimeLayout = "2006_01_02_15_04_05"

// PlaysTv _
type PlaysTv struct {
}

// NewPlaysTv _
func NewPlaysTv() *PlaysTv {
	return &PlaysTv{}
}

// IsMatch _
func (gf *PlaysTv) IsMatch(file files.Filer) bool {
	return len(PlaysTvRegexp.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *PlaysTv) Extract(file files.Filer) (time.Time, error) {
	strs := PlaysTvRegexp.FindAllString(file.Name(), -1)

	if len(strs) > 0 {
		parsedTime, err := time.ParseInLocation(
			PlaysTvTimeLayout,
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
func (gf *PlaysTv) HandlerName() string {
	return "plays_tv"
}
