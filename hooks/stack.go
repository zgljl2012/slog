package hooks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// StackHook stack hook
type StackHook struct {
	formatter   *logrus.TextFormatter
	path        string
	writeToFile bool
	lock        sync.Mutex
}

// NewStackHook new a stackHook
func NewStackHook(formatter *logrus.TextFormatter) *StackHook {
	return &StackHook{
		formatter: formatter,
	}
}

// SetLogPath write logs to the file
func (hook *StackHook) SetLogPath(path string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.path = path
	hook.writeToFile = true
}

// Levels provides the levels, only error level will print stack
func (hook *StackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called by logrus
func (hook *StackHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if entry.Level > logrus.ErrorLevel {
		if hook.writeToFile {
			hook.fileWrite(entry)
		}
		return nil
	}
	lines := strings.Split(string(debug.Stack()), "\n")
	idx := 0
	// find the first non-logrus pair
	for idx = 0; idx < (len(lines)-1)/2; idx++ {
		s := lines[1+2*idx]
		if matched, _ := regexp.MatchString("zgljl2012/slog\\.|logrus\\.|hooks\\.\\(\\*StackHook\\)\\.Fire|debug\\.Stack", s); !matched {
			break
		}
	}
	lines = append(lines[:1], lines[idx*2+1:]...)
	lines = append(lines, "Error: "+entry.Message+"\n")
	output := strings.Join(lines, "\n")
	// TODO print to the handler, not stderr
	_, _ = fmt.Fprintln(os.Stderr, output)
	if hook.writeToFile {
		hook.fileWriteRaw(output)
		hook.fileWrite(entry)
	}
	return nil
}

// Write a log line directly to a file.
func (hook *StackHook) fileWriteRaw(msg string) error {
	var (
		fd   *os.File
		path string
		err  error
	)

	path = hook.path

	dir := filepath.Dir(path)
	os.MkdirAll(dir, os.ModePerm)

	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("failed to open logfile:", path, err)
		return err
	}
	defer fd.Close()

	fd.Write([]byte(msg))
	return nil
}

// Write a log line directly to a file.
func (hook *StackHook) fileWrite(entry *logrus.Entry) error {
	var (
		msg []byte
	)

	// use our formatter instead of entry.String()
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}

	return hook.fileWriteRaw(string(msg))
}
