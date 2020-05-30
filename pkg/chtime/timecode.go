package chtime

import (
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/pkg/chtime/handlers"
	"github.com/wailorman/ffchunker/pkg/files"
	"github.com/wailorman/ffchunker/pkg/media"
)

// ErrNoTimeMatches _
var ErrNoTimeMatches = errors.New("No time matches")

// ExtractTimeHandler _
type ExtractTimeHandler interface {
	IsMatch(file files.Filer) bool
	Extract(file files.Filer) (time.Time, error)
	HandlerName() string
}

// ExtractTime _
func ExtractTime(file files.Filer) (time.Time, string, error) {
	mediaInfoGetter := media.NewInfoGetter()

	patterns := []ExtractTimeHandler{
		handlers.NewGeforceDVR(media.NewDurationCalculator(mediaInfoGetter)),
		handlers.NewGeforceFull(),
		handlers.NewAverMedia(),
		handlers.NewPlaysTv(),
		handlers.NewAction4(),
	}

	for _, pattern := range patterns {
		if pattern.IsMatch(file) {
			parsedTime, err := pattern.Extract(file)

			if err == nil {
				return parsedTime, pattern.HandlerName(), nil
			}
		}
	}

	return time.Time{}, "", ErrNoTimeMatches
}
