package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

type TestCase struct {
	name           string
	input          string
	expectedOutput *pb.ConvertTaskProgress
}

var table = []TestCase{
	{
		"first",
		"frame=   67 fps= 14 q=-1.0 Lsize=    3028kB time=00:00:01.23 bitrate=20049.4kbits/s speed=0.257x",
		&pb.ConvertTaskProgress{
			Frame:   67,
			Fps:     14.0,
			Time:    1230,
			Bitrate: (20049.4 * 1000),
			Speed:   0.257,
		},
	},
}

func Test__parseLogLine(tg *testing.T) {
	for _, testCase := range table {
		tg.Run(testCase.name, func(t *testing.T) {
			result := parseLogLine(testCase.input)
			assert.Equal(t, testCase.expectedOutput, result)
		})
	}
}
