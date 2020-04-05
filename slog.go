package slog

import (
	"os"

	"github.com/zgljl2012/slog/hook"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log instance
var log = logrus.New()
var customFormatter *logrus.TextFormatter

func init() {
	customFormatter = new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	// Add this line for logging filename and line number!
	log.SetReportCaller(false)
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	log.AddHook(hook.NewStackHook())
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
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: path,
		logrus.InfoLevel:  path,
		logrus.WarnLevel:  path,
		logrus.ErrorLevel: path,
		logrus.FatalLevel: path,
	}
	log.AddHook(lfshook.NewHook(
		pathMap,
		customFormatter,
	))
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
