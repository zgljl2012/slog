package log_test

import (
	log "slog/slog"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {

	log.SetLevel(logrus.DebugLevel)

	log.Info("Hello")
	log.Info("Hello", "this should be a error")
	log.Info("Hello", "status", "normal")

	log.Debug("debug")
	log.Debug("debug", "level", 0)

	log.Info("info")
	log.Info("info", "l", 1, "hello", "world")

	log.Warn("warn")
	log.Warn("warn", "l", 2, "hello", "world")

	log.Error("error")
	log.Error("error", "l", 3)

	t.Skip()
	log.Fatal("fatal")
	log.Fatal("fatal", "l", 4)
}
