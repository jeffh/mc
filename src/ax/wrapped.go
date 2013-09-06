package ax

import (
	"fmt"
	"sync"
	"time"
)

// Basic embeddable Logger type that can accept another logger to wrap
type BasicWrappedLogger struct {
	Logger Logger
}

func (l *BasicWrappedLogger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(format, v...)
}

func (l *BasicWrappedLogger) SetLogger(logger Logger) {
	l.Logger = logger
}

func (l *BasicWrappedLogger) WrappedLogger() Logger {
	return l.Logger
}

// A WrapLogger that ensures only one message is being logged to the
// wrapped logger at a time. This is done by locking before each call
// to the wrapped logger
type LockedLogger struct {
	BasicWrappedLogger
	mutex sync.Mutex
}

func NewLockedLogger() *LockedLogger {
	return &LockedLogger{}
}

func (l *LockedLogger) Printf(format string, v ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Logger.Printf(format, v...)
}

// A Logger that prefixes the log messages it has been given
// with a given string.
type PrefixLogger struct {
	BasicWrappedLogger
	Prefix string
}

func NewPrefixLogger(prefix string) *PrefixLogger {
	return &PrefixLogger{Prefix: prefix}
}

func (l *PrefixLogger) Printf(format string, v ...interface{}) {
	output := fmt.Sprintf(format, v...)
	l.Logger.Printf("%s%s", l.Prefix, output)
}

// A WrapLogger that logs messages with the current time prefixed
// in front of it
type TimestampLogger struct {
	BasicWrappedLogger
	Now func() string
}

func NowAsString() string {
	return time.Now().Format(time.RFC3339) + " "
}

func NewTimestampLogger() *TimestampLogger {
	return &TimestampLogger{Now: NowAsString}
}

func (l *TimestampLogger) Printf(format string, v ...interface{}) {
	output := fmt.Sprintf(format, v...)
	l.Logger.Printf("%s%s", NowAsString(), output)
}
