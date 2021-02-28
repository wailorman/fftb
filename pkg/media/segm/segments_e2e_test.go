package segm_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/segm"
	"golang.org/x/sync/errgroup"
)

func Test__splitAndConcat(t *testing.T) {
	testTable := []struct {
		file files.Filer
	}{
		{file: files.NewFile("/Users/wailorman/projects/fftb/tmp/video/video01.mp4")},
	}

	for _, testItem := range testTable {
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			segmentsDir := files.NewTempPath("fftb_test_segments")

			t.Log("segments path:", segmentsDir.FullPath())

			sliceOperation := segm.NewSliceOperation(ctx)

			err := sliceOperation.Init(segm.SliceRequest{
				InFile:         testItem.file,
				OutPath:        segmentsDir,
				KeepTimestamps: true,
				SegmentSec:     2,
			})

			assert.Nil(t, err)

			sProgress, sSegments, sFailures := sliceOperation.Run()

			segments := make([]*segm.Segment, 0)

			sliceGroup := new(errgroup.Group)
			sliceGroup.Go(func() error {
				for {
					select {
					case p := <-sProgress:
						if p != nil {
							t.Log("Slicing progress:", p.Progress())
						}

					case segment := <-sSegments:
						if segment != nil {
							segments = append(segments, segment)
						}

					case failure, failed := <-sFailures:
						if !failed {
							return nil
						}

						return failure

					// TODO: rewrite with context
					case <-time.After(5 * time.Minute):
						return errors.New("timeout reached")
					}
				}
			})

			err = sliceGroup.Wait()

			assert.Nil(t, err)

			assert.Equal(t, 69, len(segments), "Segments count")

			concatOperation := segm.NewConcatOperation(ctx)
			concatOutputFile, err := files.NewTempFile("fftb_concat_test", "fftb_concat_output.mp4")
			concatOutputFile.Create()

			assert.Nil(t, err)

			t.Log("concatenation result file path:", concatOutputFile.FullPath())

			err = concatOperation.Init(segm.ConcatRequest{
				OutFile:  concatOutputFile,
				Segments: segments,
			})

			assert.Nil(t, err)

			cProgress, cFailures := concatOperation.Run()

			concatGroup := new(errgroup.Group)
			concatGroup.Go(func() error {
				for {
					select {
					case p := <-cProgress:
						if p != nil {
							t.Logf("Concatenation progress: %f", p.Progress())
						}

					case failure, failed := <-cFailures:
						if !failed {
							return nil
						}

						return failure

					// TODO: rewrite with context
					case <-time.After(5 * time.Minute):
						return errors.New("timeout reached")
					}
				}
			})

			err = concatGroup.Wait()

			time.Sleep(5 * time.Second)

			assert.Nil(t, err)

			inputSize, err := testItem.file.Size()
			assert.Nil(t, err)

			outputSize, err := concatOutputFile.Size()
			assert.Nil(t, err)

			assert.Equal(t, int(inputSize/1024/1024), int(outputSize/1024/1024))

			segmentsDir.Destroy()
			concatOutputFile.Remove()
			sliceOperation.Purge()
			concatOperation.Prune()
		}()
	}
}
