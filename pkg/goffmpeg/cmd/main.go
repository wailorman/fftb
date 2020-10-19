package main

import (
	"fmt"

	"github.com/wailorman/fftb/pkg/goffmpeg/transcoder"
)

func main() {
	// Create new instance of transcoder
	trans := new(transcoder.Transcoder)

	// Initialize transcoder passing the input file path and output file path
	err := trans.Initialize(
		"/Users/wailorman/Desktop/W Action 07-04-2019 17-44-54.mp4",
		"/Users/wailorman/Desktop/W Action 07-04-2019 17-44-54_vtb.mp4",
	)

	if err != nil {
		panic(err)
	}

	trans.MediaFile().SetVsync(true)
	trans.MediaFile().SetHardwareAcceleration("videotoolbox")
	trans.MediaFile().SetVideoBitRate("20M")
	trans.MediaFile().SetVideoCodec("h264_videotoolbox")

	// Start transcoder process to check progress
	done := trans.Run(true)

	// Returns a channel to get the transcoding progress
	progress := trans.Output()

	for {
		select {
		case progressMessage := <-progress:
			fmt.Printf("progressMessage: %#v\n", progressMessage)
		case err := <-done:
			panic(err)
		}
	}
}
