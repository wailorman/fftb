package dlog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"
	"github.com/wailorman/fftb/pkg/ctxlog"
)

const PrefixSegmConcatOperation = "fftb.segm.concat_operation"
const PrefixSegmSliceOperation = "fftb.segm.slice_operation"
const PrefixRun = "run"
const PrefixRclone = "rclone"
const PrefixFFmpeg = "ffmpeg"
const PrefixTwirp = "twirp"

const KeyProgress = "progress"
const KeySpeed = "speed"
const KeyFPS = "fps"
const KeyCallee = "callee"
const KeyStdLevel = "stdlevel"
const KeyCommand = "command"
const KeyThread = "thread"
const KeyTaskID = "task_id"
const KeyRunID = "run_id"
const KeyPath = "path"
const KeyType = "type"

const KeyRemoteDir = "remote_dir"
const KeyLocalPath = "local_path"

const StdLevelStdout = "stdout"
const StdLevelStderr = "stderr"

const CalleeRclone = "rclone"
const CalleeFFmpeg = "ffmpeg"

func JSON(v interface{}) string {
	if v == nil {
		return "{}"
	}

	bytes, err := json.Marshal(v)

	if err != nil {
		return fmt.Sprintf("<failed to generate json: %s>", err)
	}

	return string(bytes)
}

func TwirpLogInterceptor(l *logrus.Entry) twirp.Interceptor {
	defaultLogger := ctxlog.WithPrefix(l, PrefixTwirp)

	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			serviceName, _ := twirp.ServiceName(ctx)
			methodName, _ := twirp.MethodName(ctx)

			logger := ctxlog.FromContext(ctx, PrefixTwirp)

			if logger == nil {
				logger = defaultLogger
			}

			logger.WithFields(logrus.Fields{
				"service": serviceName,
				"method":  methodName,
				"request": JSON(req),
			}).Trace("twirp request")

			resp, err := next(ctx, req)

			logger.WithFields(logrus.Fields{
				"response": JSON(resp),
				"error":    err,
			}).Trace("twirp response")

			return resp, err
		}
	}
}
