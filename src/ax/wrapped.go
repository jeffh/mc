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

// A Logger that applies a function Map to its incoming message
// before sending to its wrapped logger
type MapLogger struct {
	BasicWrappedLogger
	Map func(format string, v ...interface{}) string
}

func NewMapLogger(m func(format string, v ...interface{}) string) *MapLogger {
	return &MapLogger{Map: m}
}

func (l *MapLogger) Printf(format string, v ...interface{}) {
	l.Logger.Printf("%s", l.Map(format, v...))
}

// A Logger that prefixes the log messages it has been given
// with a given string.
func NewPrefixLogger(prefix string) *MapLogger {
	return NewMapLogger(func(format string, v ...interface{}) string {
		return prefix + fmt.Sprintf(format, v...)
	})
}

// A Logger that prefixes the current time to the log message
func NewTimestampLogger() *MapLogger {
	return NewMapLogger(func(format string, v ...interface{}) string {
		return time.Now().Format(time.RFC3339) + " " + fmt.Sprintf(format, v...)
	})
}
