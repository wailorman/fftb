package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test__GetIsLocked__Locked__NotExpired(t *testing.T) {
	until := time.Now().Add(10 * time.Hour)

	segment := &ConvertSegment{
		LockedBy:    &Author{Name: "test"},
		LockedUntil: &until,
	}

	assert.Equal(t, true, segment.GetIsLocked())
}

func Test__GetIsLocked__Locked__Expired(t *testing.T) {
	until := time.Now().Add(-(10 * time.Hour))

	segment := &ConvertSegment{
		LockedBy:    &Author{Name: "test"},
		LockedUntil: &until,
	}

	assert.Equal(t, false, segment.GetIsLocked())
}

func Test__GetIsLocked__Locked__NeverLocked(t *testing.T) {
	segment := &ConvertSegment{}

	assert.Equal(t, false, segment.GetIsLocked())
}
