package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__AverMedia__Extract(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		filename     string
		expectedTime string
	}{
		{
			filename:     "Far Cry® New Dawn 2020.02.12 - 23.03.10.00_some_hevc.mp4",
			expectedTime: "2020-02-12T23:03:10",
		},
	}

	for _, testItem := range testTable {
		timeObj, err := NewGeforceFull().Extract(newFilerStub("", testItem.filename))

		assert.Nil(err)

		assert.Equal(testItem.expectedTime, timeObj.Format("2006-01-02T15:04:05"), testItem.filename)
	}
}
