package chtime

import (
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/chunky/pkg/files"
)

// ChTimer _
type ChTimer struct {
	file files.Filer
}

// NewChTimer _
func NewChTimer(file files.Filer) *ChTimer {
	return &ChTimer{
		file: file,
	}
}

// ChTimerResult _
type ChTimerResult struct {
	Ok          bool
	Time        time.Time
	File        files.Filer
	UsedHandler string
	Error       error
}

func newChTimerResult(ok bool, time time.Time, file files.Filer, usedHandler string, err error) ChTimerResult {
	return ChTimerResult{
		Ok:          ok,
		Time:        time,
		File:        file,
		UsedHandler: usedHandler,
		Error:       err,
	}
}

// Perform _
func (t *ChTimer) Perform() ChTimerResult {
	extractedTime, usedHandler, err := ExtractTime(t.file)

	if err != nil {
		return newChTimerResult(false, extractedTime, t.file, usedHandler, err)
	}

	err = t.file.SetChTime(extractedTime)

	if err != nil {
		return newChTimerResult(false, extractedTime, t.file, usedHandler, err)
	}

	return newChTimerResult(true, extractedTime, t.file, usedHandler, nil)
}

// RecursiveChTimer _
type RecursiveChTimer struct {
	path files.Pather
}

// NewRecursiveChTimer _
func NewRecursiveChTimer(path files.Pather) *RecursiveChTimer {
	return &RecursiveChTimer{
		path: path,
	}
}

// Perform _
func (rt *RecursiveChTimer) Perform() (chan ChTimerResult, chan bool) {
	results := make(chan ChTimerResult, 0)
	done := make(chan bool, 1)

	files, err := rt.path.Files()

	if err != nil {
		panic(errors.Wrap(err, "Getting files from path"))
	}

	go func() {
		for _, file := range files {
			results <- NewChTimer(file).Perform()
		}

		done <- true
	}()

	return results, done
}
