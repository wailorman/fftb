package ctxlog

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// New _
func New(contextName string) *logrus.Entry {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Formatter = new(prefixed.TextFormatter)

	return log.
		WithField("prefix", contextName)
}

// DefaultContext _
const DefaultContext = "chunky"
