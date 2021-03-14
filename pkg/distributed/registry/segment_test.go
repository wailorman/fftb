package registry

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/distributed/ukvs/localfile"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
)

var someTime = time.Date(2020, 02, 02, 12, 10, 00, 00, time.Local)

var segmentsTestTable = []struct{ segment models.ISegment }{
	{
		&models.ConvertSegment{
			Identity:                   "segment_id_1",
			OrderIdentity:              "order_id_1",
			Type:                       models.ConvertV1Type,
			InputStorageClaimIdentity:  "local/some_dir/some_in_file",
			OutputStorageClaimIdentity: "local/some_dir/some_out_file",
			State:                      models.SegmentStatePublished,

			Publisher: &models.Author{Name: "somePublisher"},

			Params: convert.Params{
				HWAccel:          "nvenc",
				KeyframeInterval: 30,
				Preset:           "slow",
				Scale:            "1/2",
				VideoBitRate:     "10M",
				VideoQuality:     30,
				VideoCodec:       "h264",
			},
			Muxer:    "mp4",
			Position: 5,

			LockedUntil: &someTime,
			LockedBy:    &models.Author{Name: "v1/publishers/local/0009"},

			RetriesCount: 5,
			RetryAt:      &someTime,
		},
	},
}

func Test__Segment__Marshaling(t *testing.T) {
	for i, testItem := range segmentsTestTable {
		originalSegment := testItem.segment

		segmentBytes, err := marshalSegmentModel(originalSegment)

		assert.Nilf(t, err, "item %d: marshalSegmentModel error", i)

		newSegment, err := unmarshalSegmentModel(segmentBytes)

		assert.Nilf(t, err, "item %d: unmarshalSegmentModel error", i)

		assert.Equalf(t, "", cmp.Diff(originalSegment, newSegment), "item %d: diff", i)
	}
}

func Test__Segment__Persisting(t *testing.T) {
	tmpPath, err := files.NewTempFile("fftb", "test__segment__persisting.json")
	assert.Nil(t, err, "init store file")
	store, err := localfile.NewClient(context.TODO(), tmpPath.FullPath())
	assert.Nil(t, err, "localfile initialization")

	registry, err := NewRegistry(context.TODO(), store)

	assert.Nil(t, err, "registry initialization")

	for i, testItem := range segmentsTestTable {
		originalSegment := testItem.segment

		err := registry.PersistSegment(originalSegment)

		assert.Nilf(t, err, "item %d: registry.PersistSegment error", i)

		newSegment, err := registry.FindSegmentByID(originalSegment.GetID())

		assert.Nilf(t, err, "item %d: registry.FindSegmentByID error", i)
		assert.Equalf(t, "", cmp.Diff(originalSegment, newSegment), "item %d: diff", i)
	}
}
