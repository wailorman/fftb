package ffchunker

import (
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/files"
	"github.com/wailorman/ffchunker/handlers"
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
	patterns := []ExtractTimeHandler{
		handlers.NewGeforceDVR(NewDurationCalculator()),
		handlers.NewGeforceFull(),
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
