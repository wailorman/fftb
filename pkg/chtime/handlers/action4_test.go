package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__Action4__Extract(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		filename     string
		expectedTime string
	}{
		{
			filename:     "NMS 22-05-2020 22-03-32.webcam.mp4",
			expectedTime: "2020-05-22T22:03:32",
		},
	}

	for _, testItem := range testTable {
		timeObj, err := NewAction4().Extract(newFilerStub("", testItem.filename))

		assert.Nil(err)

		assert.Equal(testItem.expectedTime, timeObj.Format("2006-01-02T15:04:05"), testItem.filename)
	}
}
