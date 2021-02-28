package log

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wailorman/fftb/pkg/ctxlog"
)

func getLogrusLevel(level int) logrus.Level {
	switch level {
	case 0:
		return logrus.PanicLevel
	case 1:
		return logrus.FatalLevel
	case 2:
		return logrus.ErrorLevel
	case 3:
		return logrus.WarnLevel
	case 4:
		return logrus.InfoLevel
	case 5:
		return logrus.DebugLevel
	case 6:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

// SetLoggingLevel _
func SetLoggingLevel(c *cli.Context) {
	lvl := getLogrusLevel(c.Int("verbosity"))

	ctxlog.SetLevel(lvl)
}
