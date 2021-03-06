package slog

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zgljl2012/slog/hooks"
)

// Log instance
var log = logrus.New()
var customFormatter *logrus.TextFormatter
var hook *hooks.StackHook

func init() {
	customFormatter = new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	// Add this line for logging filename and line number!
	log.SetReportCaller(false)
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	hook = hooks.NewStackHook(customFormatter)

	log.AddHook(hook)
}

func exec(level logrus.Level, msg interface{}, args ...interface{}) {
	if args != nil {
		if len(args)%2 == 1 {
			log.WithField("msg", msg).Error("Args'count of this log should be even")
			return
		}
		var tmp *logrus.Entry
		for i := 0; i < len(args); {
			if tmp == nil {
				tmp = log.WithField((args[i]).(string), args[i+1])
			} else {
				tmp = tmp.WithField((args[i]).(string), args[i+1])
			}
			i += 2
		}
		if level == logrus.InfoLevel {
			tmp.Info(msg)
		} else if level == logrus.DebugLevel {
			tmp.Debug(msg)
		} else if level == logrus.WarnLevel {
			tmp.Warn(msg)
		} else if level == logrus.ErrorLevel {
			tmp.Error(msg)
		} else if level == logrus.FatalLevel {
			tmp.Fatal(msg)
		}
	} else {
		if level == logrus.InfoLevel {
			log.Info(msg)
		} else if level == logrus.DebugLevel {
			log.Debug(msg)
		} else if level == logrus.WarnLevel {
			log.Warn(msg)
		} else if level == logrus.ErrorLevel {
			log.Error(msg)
		} else if level == logrus.FatalLevel {
			log.Fatal(msg)
		}
	}
}

// SetLevel set log level
func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

// SetLogPath set path
func SetLogPath(path string) {
	hook.SetLogPath(path)
}

// SetRotationTime set rotation time, at least 1 second
func SetRotationTime(rotationTime time.Duration) error {
	if rotationTime < time.Second {
		return fmt.Errorf("rotation can't be less than one second")
	}
	hook.SetRotationTime(rotationTime)
	return nil
}

// Debug level
func Debug(msg interface{}, args ...interface{}) {
	exec(logrus.DebugLevel, msg, args...)
}

// Info level
func Info(msg interface{}, args ...interface{}) {
	exec(logrus.InfoLevel, msg, args...)
}

// Warn level
func Warn(msg interface{}, args ...interface{}) {
	exec(logrus.WarnLevel, msg, args...)
}

// Error level
func Error(msg interface{}, args ...interface{}) {
	exec(logrus.ErrorLevel, msg, args...)
}

// Fatal level
func Fatal(msg interface{}, args ...interface{}) {
	exec(logrus.FatalLevel, msg, args...)
}
