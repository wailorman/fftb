package handlers

import (
	"regexp"
	"time"

	"github.com/wailorman/chunky/pkg/files"
)

// Action4Regexp _
var Action4Regexp = regexp.MustCompile("(\\d{2}-\\d{2}-\\d{4} \\d{2}-\\d{2}-\\d{2})")

// NMS 22-05-2020 21-52-13.mp4
// NMS 22-05-2020 21-52-13.m4a
// NMS 22-05-2020 22-03-32.webcam.mp4

// Action4TimeLayout _
const Action4TimeLayout = "02-01-2006 15-04-05"

// Action4 _
type Action4 struct {
}

// NewAction4 _
func NewAction4() *Action4 {
	return &Action4{}
}

// IsMatch _
func (gf *Action4) IsMatch(file files.Filer) bool {
	return len(Action4Regexp.FindAllString(file.Name(), -1)) > 0
}

// Extract _
func (gf *Action4) Extract(file files.Filer) (time.Time, error) {
	strs := Action4Regexp.FindAllString(file.Name(), -1)

	if len(strs) > 0 {
		parsedTime, err := time.ParseInLocation(
			Action4TimeLayout,
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
func (gf *Action4) HandlerName() string {
	return "action4"
}
