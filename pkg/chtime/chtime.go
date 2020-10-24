package chtime

import (
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

// Instance _
type Instance struct {
	file files.Filer
}

// New _
func New(file files.Filer) *Instance {
	return &Instance{
		file: file,
	}
}

// Result _
type Result struct {
	Ok          bool
	Time        time.Time
	File        files.Filer
	UsedHandler string
	Error       error
}

func newResult(ok bool, time time.Time, file files.Filer, usedHandler string, err error) Result {
	return Result{
		Ok:          ok,
		Time:        time,
		File:        file,
		UsedHandler: usedHandler,
		Error:       err,
	}
}

// Perform _
func (t *Instance) Perform() Result {
	extractedTime, usedHandler, err := ExtractTime(t.file)

	if err != nil {
		return newResult(false, extractedTime, t.file, usedHandler, err)
	}

	err = t.file.SetChTime(extractedTime)

	if err != nil {
		return newResult(false, extractedTime, t.file, usedHandler, err)
	}

	return newResult(true, extractedTime, t.file, usedHandler, nil)
}

// RecursiveInstance _
type RecursiveInstance struct {
	path files.Pather
}

// NewRecursive _
func NewRecursive(path files.Pather) *RecursiveInstance {
	return &RecursiveInstance{
		path: path,
	}
}

// Perform _
func (rt *RecursiveInstance) Perform() (chan Result, chan bool) {
	results := make(chan Result, 0)
	done := make(chan bool, 1)

	files, err := rt.path.Files()

	if err != nil {
		panic(errors.Wrap(err, "Getting files from path"))
	}

	go func() {
		for _, file := range files {
			results <- New(file).Perform()
		}

		done <- true
	}()

	return results, done
}
