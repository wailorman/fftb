package ffmpeg

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
	"github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/run"
)

type FFmpegClient struct {
	logger      *logrus.Entry
	workingDir  string
	ffmpegPath  string
	ffprobePath string
}

var defaultFfmpegOpts = []string{"-y"}

var progressBitrateMultipliers = map[string]float64{
	"mbits/s": 1_000_000.0,
	"kbits/s": 1_000.0,
	"bits/s":  1.0,
}

func NewFFmpegClient() *FFmpegClient {
	return &FFmpegClient{
		logger:      ctxlog.New(dlog.PrefixFFmpeg).WithField(dlog.KeyCallee, dlog.CalleeFFmpeg),
		ffmpegPath:  "ffmpeg",
		ffprobePath: "ffprobe",
	}
}

func (fc *FFmpegClient) SetLogger(logger *logrus.Entry) {
	fc.logger = ctxlog.WithPrefix(logger, dlog.PrefixFFmpeg).WithField(dlog.KeyCallee, dlog.CalleeFFmpeg)
}

func (fc *FFmpegClient) SetWorkingDir(dir string) {
	fc.workingDir = dir
}

func (fc *FFmpegClient) SetFFmpegPath(path string) {
	fc.ffmpegPath = path
}

func (fc *FFmpegClient) SetFFprobePath(path string) {
	fc.ffprobePath = path
}

func (fc *FFmpegClient) GetMetadata(ctx context.Context, filePath string) (models.Metadata, error) {
	command := []string{fc.ffprobePath, "-i", filePath, "-print_format", "json", "-show_format", "-show_streams", "-show_error"}
	operation := run.New(command)
	operation.SetLogger(fc.logger)

	if fc.workingDir != "" {
		operation.SetWorkingDir(fc.workingDir)
	}

	err := operation.Run(ctx)

	if err != nil {
		return models.Metadata{}, err
	}

	stdoutLines, _, err := operation.WaitOutput()

	if err != nil {
		return models.Metadata{}, err
	}

	stdout := strings.Join(stdoutLines, "\n")

	metadata := &models.Metadata{}
	if err = json.Unmarshal([]byte(stdout), metadata); err != nil {
		return models.Metadata{}, err
	}

	return *metadata, nil
}

func (fc *FFmpegClient) Transcode(ctx context.Context, opts []string, progress chan *pb.ConvertTaskProgress) error {
	if progress != nil {
		defer close(progress)
	}

	command := []string{fc.ffmpegPath}
	command = append(command, defaultFfmpegOpts...)
	command = append(command, opts...)

	operation := run.New(command)
	operation.SetLogger(fc.logger)

	if fc.workingDir != "" {
		operation.SetWorkingDir(fc.workingDir)
	}

	err := operation.Run(ctx)

	if err != nil {
		return err
	}

	doneCh, stdoutCh, stderrCh, failuresCh := operation.StreamOutput()
	failures := make([]error, 0)

	for {
		select {
		case line, ok := <-stdoutCh:
			if ok {
				pMsg := parseLogLine(line)

				if pMsg != nil && progress != nil {
					progress <- pMsg
				}
			}

		case line, ok := <-stderrCh:
			if ok {
				pMsg := parseLogLine(line)

				if pMsg != nil && progress != nil {
					progress <- pMsg
				}
			}

		case errorLine, ok := <-failuresCh:
			if ok {
				failures = append(failures, errorLine)
			}

		case <-doneCh:
			return ctxlog.ConcatErrors(failures)
		}
	}
}

func parseLogLine(line string) *pb.ConvertTaskProgress {
	if strings.Index(line, "frame=") != 0 {
		return nil
	}

	eqRegexp := regexp.MustCompile(`=\s+`)
	strippedLine := eqRegexp.ReplaceAllString(line, "=")

	frameMatches := regexp.MustCompile(`frame=([^\s]+)`).FindAllString(strippedLine, 1)
	fpsMatches := regexp.MustCompile(`fps=([^\s]+)`).FindAllString(strippedLine, 1)
	timeMatches := regexp.MustCompile(`time=([^\s]+)`).FindAllString(strippedLine, 1)
	bitrateMatches := regexp.MustCompile(`bitrate=([^\s]+)`).FindAllString(strippedLine, 1)
	speedMatches := regexp.MustCompile(`speed=([^\s]+)`).FindAllString(strippedLine, 1)

	noEntries := len(frameMatches) == 0 && len(fpsMatches) == 0 &&
		len(timeMatches) == 0 && len(bitrateMatches) == 0 && len(speedMatches) == 0

	if noEntries {
		return nil
	}

	progressMessage := &pb.ConvertTaskProgress{}

	if len(frameMatches) > 0 {
		progressMessage.Frame, _ = strconv.ParseInt(extractDigits(frameMatches[0]), 10, 64)
	}

	if len(fpsMatches) > 0 {
		progressMessage.Fps, _ = strconv.ParseFloat(extractDigits(fpsMatches[0]), 64)
	}

	if len(timeMatches) > 0 {
		timeSegments := regexp.
			MustCompile(`(?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2}.?\d{0,4})`).
			FindAllStringSubmatch(timeMatches[0], 4)

		if timeSegments != nil {
			seconds, _ := strconv.ParseFloat(timeSegments[0][3], 64)
			minutes, _ := strconv.Atoi(timeSegments[0][2])
			hours, _ := strconv.Atoi(timeSegments[0][1])

			progressMessage.Time =
				int64(seconds*1000) +
					int64(minutes*60_000) + int64(hours*60*60_000)
		}
	}

	if len(bitrateMatches) > 0 {
		for suffix, multiplier := range progressBitrateMultipliers {
			if strings.Contains(bitrateMatches[0], suffix) {
				floatBitrate, _ := strconv.ParseFloat(extractDigits(bitrateMatches[0]), 64)
				progressMessage.Bitrate = floatBitrate * multiplier
				break
			}
		}
	}

	if len(speedMatches) > 0 {
		progressMessage.Speed, _ = strconv.ParseFloat(extractDigits(speedMatches[0]), 64)
	}

	return progressMessage
}

func extractDigits(s string) string {
	return regexp.MustCompile(`\d+.?\d+`).FindString(s)
}
