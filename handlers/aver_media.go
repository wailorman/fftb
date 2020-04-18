package handlers

import (
	"regexp"
	"time"

	"github.com/wailorman/ffchunker/files"
)

// AverMediaRegexp _
var AverMediaRegexp = regexp.MustCompile("(\\d{8}\\_\\d{6})")

// 20180505_170735.mp4
// 20060102_150405

// AverMediaTimeLayout _
const AverMediaTimeLayout = "20060102_150405"

// AverMedia _
type AverMedia struct {
}

// NewAverMedia _
func NewAverMedia() *AverMedia {
	return &AverMedia{}
}

// IsMatch _
func (gf *AverMedia) IsMatch(file files.Filer) bool {
	return len(AverMediaRegexp.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *AverMedia) Extract(file files.Filer) (time.Time, error) {
	strs := AverMediaRegexp.FindAllString(file.Name(), -1)

	if len(strs) > 0 {
		parsedTime, err := time.ParseInLocation(
			AverMediaTimeLayout,
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
func (gf *AverMedia) HandlerName() string {
	return "aver_media"
}
