package convert_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/convert"
	"github.com/wailorman/fftb/pkg/media/minfo"
	"golang.org/x/sync/errgroup"
)

func Test__batchConvert(t *testing.T) {
	testTable := []struct {
		task convert.BatchTask
	}{
		{
			task: convert.BatchTask{
				Parallelism: 1,
				Tasks: []convert.Task{
					{
						InFile:       files.NewFile("/Users/wailorman/projects/fftb/tmp/video/video02.mp4"),
						OutFile:      files.NewFile("/Users/wailorman/projects/fftb/tmp/video/video02_out.mp4"),
						HWAccel:      "",
						VideoCodec:   "h264",
						Preset:       "ultrafast",
						VideoBitRate: "",
						VideoQuality: 45,
						Scale:        "",
					},
				},
			},
		},
	}

	for _, testItem := range testTable {
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mediaInfoGetter := minfo.NewGetter()

			converter := convert.NewBatchConverter(ctx, mediaInfoGetter)
			cProgress, cFailures := converter.Convert(testItem.task)

			cg := new(errgroup.Group)
			cg.Go(func() error {
				for {
					select {
					case p, ok := <-cProgress:
						if ok {
							t.Log("Converting progress:", p.Progress.Progress())
						}

					case failure, failed := <-cFailures:
						if !failed {
							return nil
						}

						return failure.Err

					case <-time.After(5 * time.Minute):
						return errors.New("timeout reached")
					}
				}
			})

			err := cg.Wait()

			assert.Nil(t, err)

			outputSize, err := testItem.task.Tasks[0].OutFile.Size()

			assert.Nil(t, err)

			assert.GreaterOrEqual(t, outputSize, 1_000)
		}()
	}
}
