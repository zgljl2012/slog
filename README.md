# slog

simple logger: a wrapper of logrus

## Usage

The first param is the message which you want to output, the rest params should be key-value pairs. This is a simple way to realize the `.WithField(key, value)` of logrus.

+ supported write logs to file
+ supported log rotation

like below:

```golang

package log_test

import (
    log "slog/slog"
    "testing"

    "github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {

    log.SetLevel(logrus.DebugLevel)

    log.Info("Hello")
    log.Info("Hello", "field", "field value")

    log.Info("Hello", "field1", "field1 value", "field2", 2)

    // set path
    log.SetLogPath("/tmp/test.log")

    // set rotationTime
    log.SetRotationTime(time.Second)
}


```
