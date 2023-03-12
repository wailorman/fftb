package run

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
)

var windowsCancellationErrorRegexp = regexp.MustCompile(`0xc000013a`)

type Instance struct {
	command        []string
	alreadyStarted bool
	workingDir     string

	stdOutPipe io.ReadCloser
	stdErrPipe io.ReadCloser
	stdInPipe  io.WriteCloser

	done     chan struct{}
	stdout   chan string
	stderr   chan string
	failures chan error

	logger *logrus.Entry
	wg     chwg.WaitGrouper
}

func New(command []string) *Instance {
	return &Instance{
		command: command,
		logger:  ctxlog.New(dlog.PrefixRun),
		wg:      chwg.New(),
	}
}

var ErrDirty = errors.New("Runner was already started")

func (i *Instance) SetLogger(logger *logrus.Entry) {
	i.logger = ctxlog.WithPrefix(logger, dlog.PrefixRun)
}

func (i *Instance) SetWorkingDir(workingDir string) {
	i.workingDir = workingDir
}

func (i *Instance) StreamOutput() (done chan struct{}, stdout chan string, stderr chan string, failures chan error) {
	return i.done, i.stdout, i.stderr, i.failures
}

func (i *Instance) WaitOutput() (stdout []string, stderr []string, err error) {
	stdout = make([]string, 0)
	stderr = make([]string, 0)
	failures := make([]error, 0)

	for {
		doneCh, stdoutCh, stderrCh, failuresCh := i.StreamOutput()

		select {
		case line, ok := <-stdoutCh:
			if ok {
				stdout = append(stdout, line)
			}
		case line, ok := <-stderrCh:
			if ok {
				stderr = append(stderr, line)
			}
		case errorLine, ok := <-failuresCh:
			if ok {
				failures = append(failures, errorLine)
			}
		case <-doneCh:
			return stdout, stderr, ctxlog.ConcatErrors(failures)
		}
	}
}

func (i *Instance) Run(ctx context.Context) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if i.alreadyStarted {
		return ErrDirty
	}

	i.alreadyStarted = true

	i.logger.WithField(dlog.KeyCommand, strings.Join(i.command, " ")).
		Trace("running command")

	i.done = make(chan struct{})
	i.stdout = make(chan string)
	i.stderr = make(chan string)
	i.failures = make(chan error)

	go func() {
		var err error

		proc := exec.CommandContext(ctx, i.command[0], i.command[1:len(i.command)]...)

		if i.workingDir != "" {
			proc.Dir = i.workingDir
		}

		i.stdOutPipe, err = proc.StdoutPipe()

		if err != nil {
			reportFailure(i.logger, i.failures, errors.Wrap(err, "Failed to get stderr"))
			close(i.done)
			return
		}

		i.stdErrPipe, err = proc.StderrPipe()

		if err != nil {
			reportFailure(i.logger, i.failures, errors.Wrap(err, "stderr not available"))
			close(i.done)
			return
		}

		i.stdInPipe, err = proc.StdinPipe()

		if err != nil {
			reportFailure(i.logger, i.failures, errors.Wrap(err, "stdin not available"))
			close(i.done)
			return
		}

		err = proc.Start()

		if err != nil {
			reportFailure(i.logger, i.failures, errors.Wrap(err, "Failed to run command"))
			close(i.done)
			return
		}

		i.wg.Add(2)
		go scanLines(i.wg, i.logger.WithField(dlog.KeyStdLevel, dlog.StdLevelStdout), i.stdOutPipe, i.stdout)
		go scanLines(i.wg, i.logger.WithField(dlog.KeyStdLevel, dlog.StdLevelStderr), i.stdErrPipe, i.stderr)

		i.wg.Wait()
		err = proc.Wait()

		if err != nil {
			cancellationError := errors.Is(ctx.Err(), context.Canceled) || windowsCancellationErrorRegexp.MatchString(err.Error())

			if cancellationError {
				reportFailure(i.logger, i.failures, context.Canceled)
			} else {
				reportFailure(i.logger, i.failures, errors.Wrap(err, "Failed to finish process"))
			}
		}

		close(i.failures)
		close(i.done)
	}()

	return nil
}

type Waiter interface {
	Wait() error
}

type PipeCloser interface {
	Close() error
}

func scanLines(wg chwg.WaitGrouper, logger *logrus.Entry, pipe io.ReadCloser, out chan string) {
	defer wg.Done()
	defer close(out)

	scanner := bufio.NewScanner(pipe)

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
		line := scanner.Text()
		logger.Trace(line)
		out <- line
	}
}

func reportFailure(logger *logrus.Entry, ch chan error, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}

	logger.WithError(err).Trace("failure received")
	ch <- err
}
