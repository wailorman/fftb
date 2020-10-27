package segm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/files"
)

func Test__getSegmentFromFile(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		file            files.Filer
		expectedSegment *Segment
	}{
		{
			file:            files.NewFile("/tmp/.DS_file"),
			expectedSegment: nil,
		},
		{
			file:            files.NewFile("/tmp/fftb_out_"),
			expectedSegment: nil,
		},
		{
			file:            files.NewFile("/tmp/other_out_001"),
			expectedSegment: nil,
		},
		{
			file:            files.NewFile("/tmp/fftb_out_000"),
			expectedSegment: &Segment{Position: 0, File: files.NewFile("/tmp/fftb_out_000")},
		},
		{
			file:            files.NewFile("/tmp/fftb_out_00005"),
			expectedSegment: &Segment{Position: 5, File: files.NewFile("/tmp/fftb_out_00005")},
		},
		{
			file:            files.NewFile("/tmp/fftb_out_10005"),
			expectedSegment: &Segment{Position: 10005, File: files.NewFile("/tmp/fftb_out_10005")},
		},
	}

	for i, testItem := range testTable {
		segment := getSegmentFromFile(testItem.file)

		expectedSegment := testItem.expectedSegment

		if expectedSegment == nil {
			assert.Nil(segment, i)
		} else {
			assert.NotNil(segment, i)
			position := segment.Position
			fullPath := segment.File.FullPath()

			expectedPosition := expectedSegment.Position
			expectedFullPath := expectedSegment.File.FullPath()

			assert.Equal(expectedPosition, position, i)
			assert.Equal(expectedFullPath, fullPath, i)
		}
	}

}
