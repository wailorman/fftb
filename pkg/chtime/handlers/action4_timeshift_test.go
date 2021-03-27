package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__Action4Timeshift__Extract(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		filename     string
		expectedTime string
	}{
		{
			filename:     "TimeShift 12-02-2020 23-03-10.mp4",
			expectedTime: "2020-02-12T23:00:10",
		},
	}

	for _, testItem := range testTable {
		timeObj, err := NewAction4Timeshift(newDurationCalculatorStub(180)).Extract(newFilerStub("", testItem.filename))

		assert.Nil(err)

		assert.Equal(testItem.expectedTime, timeObj.Format("2006-01-02T15:04:05"), testItem.filename)
	}
}
