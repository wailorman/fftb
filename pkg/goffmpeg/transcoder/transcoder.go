package transcoder

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/goffmpeg/ffmpeg"
	"github.com/wailorman/fftb/pkg/goffmpeg/models"
	"github.com/wailorman/fftb/pkg/goffmpeg/utils"
	"github.com/wailorman/fftb/pkg/run"
)

// Transcoder Main struct
type Transcoder struct {
	stdErrPipe    io.ReadCloser
	stdStdinPipe  io.WriteCloser
	process       *exec.Cmd
	mediafile     *models.Mediafile
	configuration ffmpeg.Configuration
	ctx           context.Context
	logger        logrus.FieldLogger
}

// LoggingPrefix _
const LoggingPrefix = "goffmpeg"

// New _
func New(ctx context.Context) *Transcoder {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, LoggingPrefix); logger == nil {
		logger = ctxlog.New(LoggingPrefix)
	}

	return &Transcoder{
		ctx:    ctx,
		logger: logger,
	}
}

// SetProcessStderrPipe Set the STDERR pipe
func (t *Transcoder) SetProcessStderrPipe(v io.ReadCloser) {
	t.stdErrPipe = v
}

// SetProcessStdinPipe Set the STDIN pipe
func (t *Transcoder) SetProcessStdinPipe(v io.WriteCloser) {
	t.stdStdinPipe = v
}

// SetProcess Set the transcoding process
func (t *Transcoder) SetProcess(cmd *exec.Cmd) {
	t.process = cmd
}

// SetMediaFile Set the media file
func (t *Transcoder) SetMediaFile(v *models.Mediafile) {
	t.mediafile = v
}

// SetConfiguration Set the transcoding configuration
func (t *Transcoder) SetConfiguration(v ffmpeg.Configuration) {
	t.configuration = v
}

// Process Get transcoding process
func (t Transcoder) Process() *exec.Cmd {
	return t.process
}

// MediaFile Get the ttranscoding media file.
func (t Transcoder) MediaFile() *models.Mediafile {
	return t.mediafile
}

// FFmpegExec Get FFmpeg Bin path
func (t Transcoder) FFmpegExec() string {
	return t.configuration.FfmpegBin
}

// FFprobeExec Get FFprobe Bin path
func (t Transcoder) FFprobeExec() string {
	return t.configuration.FfprobeBin
}

// GetCommand Build and get command
func (t Transcoder) GetCommand() []string {
	media := t.mediafile
	rcommand := append([]string{"-y"}, media.ToStrCommand()...)
	return rcommand
}

// InitializeEmptyTranscoder initializes the fields necessary for a blank transcoder
func (t *Transcoder) InitializeEmptyTranscoder() error {
	var err error
	cfg := t.configuration
	if len(cfg.FfmpegBin) == 0 || len(cfg.FfprobeBin) == 0 {
		cfg, err = ffmpeg.Configure(t.ctx)
		if err != nil {
			return err
		}
	}
	// Set new Mediafile
	MediaFile := new(models.Mediafile)

	// Set transcoder configuration
	t.SetMediaFile(MediaFile)
	t.SetConfiguration(cfg)
	return nil
}

// SetInputPath sets the input path for transcoding
func (t *Transcoder) SetInputPath(inputPath string) error {
	if t.mediafile.InputPipe() {
		return errors.New("cannot set an input path when an input pipe command has been set")
	}
	t.mediafile.SetInputPath(inputPath)
	return nil
}

// SetOutputPath sets the output path for transcoding
func (t *Transcoder) SetOutputPath(inputPath string) error {
	if t.mediafile.OutputPipe() {
		return errors.New("cannot set an input path when an input pipe command has been set")
	}
	t.mediafile.SetOutputPath(inputPath)
	return nil
}

// CreateInputPipe creates an input pipe for the transcoding process
func (t *Transcoder) CreateInputPipe() (*io.PipeWriter, error) {
	if t.mediafile.InputPath() != "" {
		return nil, errors.New("cannot set an input pipe when an input path exists")
	}
	inputPipeReader, inputPipeWriter := io.Pipe()
	t.mediafile.SetInputPipe(true)
	t.mediafile.SetInputPipeReader(inputPipeReader)
	t.mediafile.SetInputPipeWriter(inputPipeWriter)
	return inputPipeWriter, nil
}

// CreateOutputPipe creates an output pipe for the transcoding process
func (t *Transcoder) CreateOutputPipe(containerFormat string) (*io.PipeReader, error) {
	if t.mediafile.OutputPath() != "" {
		return nil, errors.New("cannot set an output pipe when an output path exists")
	}
	t.mediafile.SetOutputFormat(containerFormat)

	t.mediafile.SetMovFlags("frag_keyframe")
	outputPipeReader, outputPipeWriter := io.Pipe()
	t.mediafile.SetOutputPipe(true)
	t.mediafile.SetOutputPipeReader(outputPipeReader)
	t.mediafile.SetOutputPipeWriter(outputPipeWriter)
	return outputPipeReader, nil
}

// Initialize Init the transcoding process
func (t *Transcoder) Initialize(inputPath string, outputPath string) error {
	var err error

	cfg := t.configuration

	if len(cfg.FfmpegBin) == 0 || len(cfg.FfprobeBin) == 0 {
		cfg, err = ffmpeg.Configure(t.ctx)
		if err != nil {
			return err
		}
	}

	if inputPath == "" {
		return errors.New("error on transcoder.Initialize: inputPath missing")
	}

	// Set new Mediafile
	MediaFile := new(models.Mediafile)
	MediaFile.SetInputPath(inputPath)
	MediaFile.SetOutputPath(outputPath)

	if isFilePossiblyHasMetadata(inputPath) {
		metadata, err := t.GetFileMetadata(inputPath)

		if err != nil {
			return errors.Wrap(err, "Getting file metadata")
		}

		MediaFile.SetMetadata(metadata)
	}

	// Set transcoder configuration
	t.SetMediaFile(MediaFile)
	t.SetConfiguration(cfg)

	return nil

}

// GetFileMetadata Returns file metadata from ffprobe
func (t *Transcoder) GetFileMetadata(filePath string) (models.Metadata, error) {
	var err error
	var outb, errb bytes.Buffer
	metadata := &models.Metadata{}

	cfg := t.configuration

	if len(cfg.FfmpegBin) == 0 || len(cfg.FfprobeBin) == 0 {
		cfg, err = ffmpeg.Configure(t.ctx)
		if err != nil {
			return models.Metadata{}, err
		}
	}

	command := []string{"-i", filePath, "-print_format", "json", "-show_format", "-show_streams", "-show_error"}

	t.logger.WithField("command", fmt.Sprintf("%s %s", cfg.FfprobeBin, strings.Join(command, " "))).
		Debug("Running ffprobe")

	cmd := exec.Command(cfg.FfprobeBin, command...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()
	if err != nil {
		return models.Metadata{}, fmt.Errorf("error executing (%s %s) | error: %s | message: %s %s", cfg.FfprobeBin, command, err, outb.String(), errb.String())
	}

	if err = json.Unmarshal([]byte(outb.String()), metadata); err != nil {
		return models.Metadata{}, err
	}

	for i := range metadata.Streams {
		stream := &metadata.Streams[i]

		stream.DurationFloat, err = strconv.ParseFloat(
			stream.Duration,
			64,
		)
	}

	return *metadata, nil
}

// GetFramesMetadata _
func (t *Transcoder) GetFramesMetadata(filePath string) (chan bool, chan models.Framer, chan error) {
	done := make(chan bool)
	frames := make(chan models.Framer, 0)
	failed := make(chan error)

	go func() {
		defer close(done)
		defer close(frames)
		defer close(failed)

		command := []string{t.configuration.FfprobeBin, "-v", "quiet", "-i", filePath, "-print_format", "json=c=1", "-show_frames", "-show_format", "-show_streams"}

		t.logger.WithField("command", strings.Join(command, " ")).
			Debug("Running ffprobe")

		runner := run.New(command)
		err := runner.Run()

		if err != nil {
			failed <- errors.Wrap(err, "Starting runner")
			done <- true
			return
		}

		rDone, rStdout, rStderr, rFailures := runner.StreamOutput()

		stdErrMessages := make([]string, 0)

		for {
			select {
			case errMsg := <-rStderr:
				if errMsg != "" {
					stdErrMessages = append(stdErrMessages, errMsg)
				}
			case line := <-rStdout:
				if line != "" {
					line = strings.ReplaceAll(line, "},", "}")

					if strings.Contains(line, "\"media_type\": \"audio\"") {
						audioFrame := &models.AudioFrame{}

						err := json.Unmarshal([]byte(line), audioFrame)

						if err != nil {
							failed <- errors.Wrap(err, "Unmarshaling audio frame")
							done <- true
							return
						}

						frames <- audioFrame
					} else if strings.Contains(line, "\"media_type\": \"video\"") {
						videoFrame := &models.VideoFrame{}

						err := json.Unmarshal([]byte(line), videoFrame)

						if err != nil {
							failed <- errors.Wrap(err, "Unmarshaling video frame")
							done <- true
							return
						}

						frames <- videoFrame
					}
				}
			case err := <-rFailures:
				failed <- errors.Wrap(err, "Receiving output from ffprobe")
				done <- true
				return
			case <-rDone:
				if len(stdErrMessages) > 0 {
					errStr := strings.Join(stdErrMessages, "; ")
					failed <- errors.New(errStr)
				}

				done <- true
				return
			}

		}
	}()

	return done, frames, failed
}

// Run Starts the transcoding process
func (t *Transcoder) Run(progress bool) <-chan error {
	done := make(chan error)
	command := t.GetCommand()

	if !progress {
		command = append([]string{"-nostats", "-loglevel", "0"}, command...)
	}

	t.logger.WithField("command", fmt.Sprintf("%s %s", t.configuration.FfmpegBin, strings.Join(command, " "))).
		Debug("Running ffmpeg")

	proc := exec.Command(t.configuration.FfmpegBin, command...)
	if progress {
		errStream, err := proc.StderrPipe()
		if err != nil {
			fmt.Println("Progress not available: " + err.Error())
		} else {
			t.stdErrPipe = errStream
		}
	}

	// Set the stdinPipe in case we need to stop the transcoding
	stdin, err := proc.StdinPipe()
	if nil != err {
		fmt.Println("Stdin not available: " + err.Error())
	}

	t.stdStdinPipe = stdin

	// If the user has requested progress, we send it to them on a Buffer
	var outb, errb bytes.Buffer
	if progress {
		proc.Stdout = &outb
	}

	// If an input pipe has been set, we set it as stdin for the transcoding
	if t.mediafile.InputPipe() {
		proc.Stdin = t.mediafile.InputPipeReader()
	}

	// If an output pipe has been set, we set it as stdout for the transcoding
	if t.mediafile.OutputPipe() {
		proc.Stdout = t.mediafile.OutputPipeWriter()
	}

	err = proc.Start()

	t.SetProcess(proc)

	go func(err error) {
		if err != nil {
			done <- fmt.Errorf("Failed Start FFMPEG (%s) with %s, message %s %s", command, err, outb.String(), errb.String())
			close(done)
			return
		}

		err = proc.Wait()

		go t.closePipes()

		if err != nil {
			err = fmt.Errorf("Failed Finish FFMPEG (%s) with %s message %s %s", command, err, outb.String(), errb.String())
		}
		done <- err
		close(done)
	}(err)

	return done
}

// Stop Ends the transcoding process
func (t *Transcoder) Stop() error {
	if t.process != nil {

		stdin := t.stdStdinPipe
		if stdin != nil {
			stdin.Write([]byte("q\n"))
		}
	}
	return nil
}

// Kill kills ffmpeg process
func (t *Transcoder) Kill() error {
	if t.process != nil {
		return t.process.Process.Kill()
	}

	return nil
}

// Output Returns the transcoding progress channel
func (t Transcoder) Output() chan models.Progress {
	out := make(chan models.Progress)

	go func() {
		defer close(out)
		if t.stdErrPipe == nil {
			out <- models.Progress{}
			return
		}

		defer t.stdErrPipe.Close()

		scanner := bufio.NewScanner(t.stdErrPipe)

		split := func(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.IndexByte(data, '\n'); i >= 0 {
				// We have a full newline-terminated line.
				return i + 1, data[0:i], nil
			}
			if i := bytes.IndexByte(data, '\r'); i >= 0 {
				// We have a cr terminated line
				return i + 1, data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}

			return 0, nil, nil
		}

		scanner.Split(split)
		buf := make([]byte, 2)
		scanner.Buffer(buf, bufio.MaxScanTokenSize)

		for scanner.Scan() {
			Progress := new(models.Progress)
			line := scanner.Text()

			t.logger.WithField("line", line).
				Trace("Received line from ffmpeg output")

			if strings.Contains(line, "frame=") && strings.Contains(line, "time=") && strings.Contains(line, "bitrate=") {
				var re = regexp.MustCompile(`=\s+`)
				st := re.ReplaceAllString(line, `=`)

				f := strings.Fields(st)
				var framesProcessed string
				var currentTime string
				var currentBitrate string
				var currentSpeed string
				var fps float64

				for j := 0; j < len(f); j++ {
					field := f[j]
					fieldSplit := strings.Split(field, "=")

					if len(fieldSplit) > 1 {
						fieldname := strings.Split(field, "=")[0]
						fieldvalue := strings.Split(field, "=")[1]

						if fieldname == "frame" {
							framesProcessed = fieldvalue
						}

						if fieldname == "time" {
							currentTime = fieldvalue
						}

						if fieldname == "bitrate" {
							currentBitrate = fieldvalue
						}

						if fieldname == "speed" {
							currentSpeed = fieldvalue
						}

						if fieldname == "fps" {
							fps, _ = strconv.ParseFloat(fieldvalue, 64)
						}
					}
				}

				timesec := utils.DurToSec(currentTime)
				dursec, _ := strconv.ParseFloat(t.MediaFile().Metadata().Format.Duration, 64)
				//live stream check
				if dursec != 0 {
					// Progress calculation
					progress := (timesec * 100) / dursec
					Progress.Progress = progress
				}
				Progress.CurrentBitrate = currentBitrate
				Progress.FramesProcessed = framesProcessed
				Progress.CurrentTime = currentTime
				Progress.Speed = currentSpeed
				Progress.FPS = fps
				out <- *Progress
			}
		}
	}()

	return out
}

func (t *Transcoder) closePipes() {
	if t.mediafile.InputPipe() {
		t.mediafile.InputPipeReader().Close()
	}
	if t.mediafile.OutputPipe() {
		t.mediafile.OutputPipeWriter().Close()
	}
}

var metadataBlacklistedExtensions = []string{".txt"}

func isFilePossiblyHasMetadata(filename string) bool {
	for _, ext := range metadataBlacklistedExtensions {
		if filepath.Ext(filename) == ext {
			return false
		}
	}

	return true
}
