package minfo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

// Instance _
type Instance struct {
}

// VideoFramesSummary _
type VideoFramesSummary struct {
	IFramesCount int `json:"i_frames_count"`
	BFramesCount int `json:"b_frames_count"`
	PFramesCount int `json:"p_frames_count"`
}

// FramesSummary _
type FramesSummary struct {
	Video VideoFramesSummary `json:"video"`
}

// Getter _
type Getter interface {
	GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error)
	GetFramesSummary(file files.Filer) (FramesSummary, error)
	GetFramesList(file files.Filer) (chan bool, chan ffmpegModels.Framer, chan error)
}

// New _
func New() Getter {
	return &Instance{}
}

// GetMediaInfo _
func (ig *Instance) GetMediaInfo(file files.Filer) (ffmpegModels.Metadata, error) {
	trans := transcoder.New(context.TODO())

	err := trans.InitializeEmptyTranscoder()

	if err != nil {
		return ffmpegModels.Metadata{}, errors.Wrap(err, "Initializing ffprobe instance")
	}

	metadata, err := trans.GetFileMetadata(file.FullPath())

	if err != nil {
		return ffmpegModels.Metadata{}, errors.Wrap(err, "Getting file metadata from ffprobe")
	}

	return metadata, nil
}

// GetFramesSummary _
func (ig *Instance) GetFramesSummary(file files.Filer) (FramesSummary, error) {
	summary := &FramesSummary{}

	done, frames, failures := ig.GetFramesList(file)

	for {
		select {
		case frame := <-frames:
			switch f := frame.(type) {
			case *ffmpegModels.VideoFrame:

				switch f.PictType {
				case "I":
					summary.Video.IFramesCount++
				case "B":
					summary.Video.BFramesCount++
				case "P":
					summary.Video.PFramesCount++
				}

			}
		case failure := <-failures:
			return FramesSummary{}, errors.Wrap(failure, "Failed to get frames information")
		case <-done:
			return *summary, nil
		}
	}
}

// GetFramesList _
func (ig *Instance) GetFramesList(file files.Filer) (chan bool, chan ffmpegModels.Framer, chan error) {
	done := make(chan bool)
	frames := make(chan ffmpegModels.Framer, 0)
	failed := make(chan error, 0)

	go func() {
		defer close(done)
		defer close(frames)
		defer close(failed)

		trans := &transcoder.Transcoder{}

		err := trans.InitializeEmptyTranscoder()

		if err != nil {
			failed <- errors.Wrap(err, "Initializing ffprobe instance")
			done <- true
			return
		}

		fDone, fFrames, fFailed := trans.GetFramesMetadata(file.FullPath())

		for {
			select {
			case frame := <-fFrames:
				frames <- frame
			case failure := <-fFailed:
				failed <- failure
			case <-fDone:
				done <- true
				return
			}
		}
	}()

	return done, frames, failed
}
