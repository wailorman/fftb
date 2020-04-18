package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__GeforceFull__Extract(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		filename     string
		expectedTime string
	}{
		{
			filename:     "20180506_170735.mp4",
			expectedTime: "2018-05-06T17:07:35",
		},
	}

	for _, testItem := range testTable {
		timeObj, err := NewAverMedia().Extract(newFilerStub("", testItem.filename))

		assert.Nil(err)

		assert.Equal(testItem.expectedTime, timeObj.Format("2006-01-02T15:04:05"), testItem.filename)
	}
}
