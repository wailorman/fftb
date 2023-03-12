package stdouthook

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Hook struct {
	formatter logrus.Formatter
	level     logrus.Level
}

type HookParams struct {
	Formatter logrus.Formatter
	Level     logrus.Level
}

func New(params HookParams) *Hook {
	return &Hook{
		formatter: params.Formatter,
		level:     params.Level,
	}
}

func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	if entry.Level <= h.level {
		entryData, err := h.formatter.Format(entry)

		if err != nil {
			return err
		}

		io.WriteString(os.Stdout, string(entryData))
	}

	return nil
}
