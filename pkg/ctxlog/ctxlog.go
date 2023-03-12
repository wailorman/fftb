package ctxlog

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// New _
func New(contextName string) *logrus.Entry {
	loggerInstance.SetLevel(logrus.TraceLevel)
	loggerInstance.Formatter = new(prefixed.TextFormatter)

	return WithPrefix(loggerInstance, contextName)
}

// DefaultContext _
const DefaultContext = "fftb"

// Logger _
var Logger = New(DefaultContext)

var loggerInstance = logrus.New()

type logKey string

// LoggerContextKey _
const LoggerContextKey logKey = "logger"

// SetLevel _
func SetLevel(lvl logrus.Level) {
	loggerInstance.SetLevel(lvl)
}

// FromContext _
func FromContext(ctx context.Context, prefix string) logrus.FieldLogger {
	logger, ok := ctx.Value(LoggerContextKey).(logrus.FieldLogger)

	if !ok {
		return nil
	}

	return logger.WithField("prefix", prefix)
}

// WithPrefix _
func WithPrefix(logger logrus.FieldLogger, prefix string) *logrus.Entry {
	return logger.WithField("prefix", prefix)
}

func ConcatErrors(errs []error) error {
	if len(errs) > 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	strErrs := make([]string, 0)
	for _, err := range errs {
		strErrs = append(strErrs, err.Error())
	}

	return errors.Wrapf(errs[0], "(also received errors: `%s`)", strings.Join(strErrs, "; "))
}
