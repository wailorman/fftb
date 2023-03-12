package rclone

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/run"
)

var defaultRcloneParams = []string{
	"-v",
	"--use-json-log",
	"--stats", "1s",
}

type RcloneClient struct {
	logger     *logrus.Entry
	path       string
	configPath string
}

func NewRcloneClient() *RcloneClient {
	return &RcloneClient{
		logger: ctxlog.New(dlog.PrefixRclone),
		path:   "rclone",
	}
}

func (rc *RcloneClient) SetLogger(logger *logrus.Entry) {
	rc.logger = ctxlog.WithPrefix(logger, dlog.PrefixRclone)
}

func (rc *RcloneClient) SetConfigPath(path string) {
	rc.configPath = path
}

func (rc *RcloneClient) SetPath(path string) {
	rc.path = path
}

func (rc *RcloneClient) Pull(ctx context.Context, remotePath, localPath string, progress chan ProgressMessage) error {
	return rc.Exec(ctx, progress,
		"copy",
		remotePath,
		localPath,
	)
}

func (rc *RcloneClient) Push(ctx context.Context, localPath, remotePath string, progress chan ProgressMessage) error {
	return rc.Exec(ctx, progress,
		"copy",
		localPath,
		remotePath,
	)
}

func (rc *RcloneClient) Exec(ctx context.Context, progress chan ProgressMessage, opts ...string) error {
	command := []string{rc.path}
	command = append(command, defaultRcloneParams...)
	if rc.configPath != "" {
		command = append(command, []string{"--config", rc.configPath}...)
	}
	command = append(command, opts...)

	cmd := run.New(command)
	cmd.SetLogger(rc.logger.WithField(dlog.KeyCallee, dlog.CalleeRclone))
	err := cmd.Run(ctx)

	if err != nil {
		return errors.Wrap(err, "Running rclone command")
	}

	if progress != nil {
		defer close(progress)
	}

	doneCh, stdoutCh, stderrCh, failuresCh := cmd.StreamOutput()
	failures := make([]error, 0)
	for {
		select {
		case line, ok := <-stdoutCh:
			if ok {
				pMsg := parseLogLine(line)

				if pMsg != nil && progress != nil {
					progress <- *pMsg
				}
			}

		case line, ok := <-stderrCh:
			if ok {
				pMsg := parseLogLine(line)

				if pMsg != nil && progress != nil {
					progress <- *pMsg
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

func parseLogLine(line string) *ProgressMessage {
	firstCharacter := line[0:1]

	if firstCharacter == "{" {
		progressMessage := &ProgressMessage{}
		err := json.Unmarshal([]byte(line), progressMessage)

		if err != nil {
			panic(err)
		}

		if progressMessage.IsValid() {
			return progressMessage
		}
	}

	return nil
}
