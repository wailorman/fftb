package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__PlaysTv__Extract(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		filename     string
		expectedTime string
	}{
		{
			filename:     "2016_05_20_15_31_51-ses.mp4",
			expectedTime: "2016-05-20T15:31:51",
		},
	}

	for _, testItem := range testTable {
		timeObj, err := NewPlaysTv().Extract(newFilerStub("", testItem.filename))

		assert.Nil(err)

		assert.Equal(testItem.expectedTime, timeObj.Format("2006-01-02T15:04:05"), testItem.filename)
	}
}
