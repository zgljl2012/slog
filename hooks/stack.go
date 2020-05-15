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
	"time"

	"github.com/sirupsen/logrus"
)

// StackHook stack hook
type StackHook struct {
	formatter    *logrus.TextFormatter
	path         string
	writeToFile  bool
	lock         sync.Mutex
	rotationTime time.Duration
	lastRotate   time.Time
}

// NewStackHook new a stackHook
func NewStackHook(formatter *logrus.TextFormatter) *StackHook {
	return &StackHook{
		formatter:    formatter,
		rotationTime: -1,
		lastRotate:   time.Now(),
	}
}

// SetLogPath write logs to the file
func (hook *StackHook) SetLogPath(path string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.path = path
	hook.writeToFile = true
}

// SetRotationTime set rotation time
func (hook *StackHook) SetRotationTime(rotationTime time.Duration) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.rotationTime = rotationTime
	hook.lastRotate = time.Now()
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
	// write output to stderr
	_, _ = fmt.Fprintln(os.Stderr, output)
	if hook.writeToFile {
		hook.fileWriteRaw(output)
		hook.fileWrite(entry)
	}
	return nil
}

func (hook *StackHook) genFilename() string {
	now := time.Now()

	// XXX HACK: Truncate only happens in UTC semantics, apparently.
	// observed values for truncating given time with 86400 secs:
	//
	// before truncation: 2018/06/01 03:54:54 2018-06-01T03:18:00+09:00
	// after  truncation: 2018/06/01 03:54:54 2018-05-31T09:00:00+09:00
	//
	// This is really annoying when we want to truncate in local time
	// so we hack: we take the apparent local time in the local zone,
	// and pretend that it's in UTC. do our math, and put it back to
	// the local zone
	return hook.path + "." + now.Format("200612150405")
}

// Write a log line directly to a file.
func (hook *StackHook) fileWriteRaw(msg string) error {
	var (
		fd   *os.File
		path string
		err  error
	)

	path = hook.path

	// rotate check
	if hook.rotationTime > 0 {
		now := time.Now()
		if now.After(hook.lastRotate.Add(hook.rotationTime)) {
			// rotate now
			if err := os.Rename(path, hook.genFilename()); err != nil {
				log.Println("failed to rename logfile:", path, err)
				return err
			}
			hook.lastRotate = now
		}
	}

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
