package hook

import (
	"fmt"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/sirupsen/logrus"
)

// StackHook stack hook
type StackHook struct{}

// NewStackHook new a stackHook
func NewStackHook() StackHook {
	return StackHook{}
}

// Levels provides the levels, only error level will print stack
func (hook StackHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}

// Fire is called by logrus
func (hook StackHook) Fire(entry *logrus.Entry) error {
	lines := strings.Split(string(debug.Stack()), "\n")
	idx := 0
	// find the first non-logrus pair
	for idx = 0; idx < (len(lines)-1)/2; idx++ {
		s := lines[1+2*idx]
		if matched, _ := regexp.MatchString("logrus|commons/log|debug\\.Stack", s); !matched {
			break
		}
	}
	lines = append(lines[:1], lines[idx*2+1:]...)
	lines = append(lines, "Error: "+entry.Message+"\n")
	output := strings.Join(lines, "\n")
	// TODO print to the handler, not stderr
	_, _ = fmt.Fprintln(os.Stderr, output)
	return nil
}
