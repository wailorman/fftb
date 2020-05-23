package ffchunker

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/wailorman/ffchunker/ctxlog"
	"github.com/wailorman/ffchunker/files"
)

// VideoCutterInstance _
type VideoCutterInstance struct {
}

// NewVideoCutter _
func NewVideoCutter() *VideoCutterInstance {
	return &VideoCutterInstance{}
}

// VideoCutter _
type VideoCutter interface {
	CutVideo(inFile files.Filer, outFile files.Filer, offset float64, maxFileSize int) (files.Filer, error)
}

// CutVideo _
func (ci *VideoCutterInstance) CutVideo(
	inFile files.Filer,
	outFile files.Filer,
	offset float64,
	maxFileSize int,
) (files.Filer, error) {

	log := ctxlog.New(ctxlog.DefaultContext + ".cutter")

	cmdStr := fmt.Sprintf(
		"ffmpeg -ss %f -i \"%s\" -fs %d -c:v copy -avoid_negative_ts make_zero -c:a copy \"%s\"",
		offset,
		inFile.FullPath(),
		maxFileSize,
		outFile.FullPath(),
	)

	// fmt.Printf("cmdStr: %#v\n", cmdStr)

	log.WithField("command", cmdStr).
		Info("Running ffmpeg command...")

	output, err := exec.Command("bash", "-c", cmdStr).Output()

	if err != nil {
		return nil, errors.Wrap(err, "ffmpeg executing problem: "+string(output))
	}

	return outFile, nil
}
