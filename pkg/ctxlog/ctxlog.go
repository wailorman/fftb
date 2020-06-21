package ctxlog

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// New _
func New(contextName string) *logrus.Entry {
	loggerInstance.SetLevel(logrus.DebugLevel)
	loggerInstance.Formatter = new(prefixed.TextFormatter)

	return loggerInstance.
		WithField("prefix", contextName)
}

// DefaultContext _
const DefaultContext = "fftb"

// Logger _
var Logger = New(DefaultContext)

var loggerInstance = logrus.New()

// SetLevel _
func SetLevel(lvl logrus.Level) {
	loggerInstance.SetLevel(lvl)
}
