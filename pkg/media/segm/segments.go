package segm

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/wailorman/fftb/pkg/files"
)

// ErrNotInitialized happened when instance wasn't initialized by Init() func
var ErrNotInitialized = errors.New("operation have not been initialized")

// ErrAlreadyInitialized happened when Init() func called twice
var ErrAlreadyInitialized = errors.New("operation was already initialized")

// ErrAlreadyStarted happened when Run() func called twice
var ErrAlreadyStarted = errors.New("operation was already started")

// Segment _
type Segment struct {
	// from 0 to inf
	Position int
	File     files.Filer
}

func createSegmentsList(segs []*Segment) string {
	sort.SliceStable(segs, func(i, j int) bool {
		return segs[i].Position < segs[j].Position
	})

	textSegs := make([]string, 0)

	for _, seg := range segs {
		textSegs = append(textSegs, fmt.Sprintf("file '%s'", seg.File.FullPath()))
	}

	list := strings.Join(textSegs, "\n")

	return list
}

func collectSegments(files []files.Filer) []*Segment {
	result := make([]*Segment, 0)

	for _, file := range files {
		foundSegment := getSegmentFromFile(file)

		if foundSegment != nil {
			result = append(result, foundSegment)
		}
	}

	return result
}

func getSegmentFromFile(file files.Filer) *Segment {
	fileName := file.Name()

	reFull := regexp.MustCompile(segmentPrefix + `\d+`)
	reNumber := regexp.MustCompile(`\d+`)

	if !reFull.MatchString(fileName) {
		return nil
	}

	foundStrNum := reNumber.FindString(fileName)

	number, err := strconv.Atoi(foundStrNum)

	if err != nil {
		return nil
	}

	return &Segment{
		Position: number,
		File:     file,
	}
}

func createTmpSubdir(dir files.Pather) (files.Pather, error) {
	id := fmt.Sprint(rand.Int())
	tmpPath := dir.BuildSubpath("_fftb_segments_" + id)
	err := tmpPath.Create()

	if err != nil {
		return nil, err
	}

	return tmpPath, nil
}
