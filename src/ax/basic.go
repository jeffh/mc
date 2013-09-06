package ax

import (
	"fmt"
)

// A Null Logger is a logger that doesn't actually log.
type NullLogger struct{}

func (l *NullLogger) Printf(format string, values ...interface{}) {
}

func (l *NullLogger) SetLogger(logger Logger) {
}

func (l *NullLogger) WrappedLogger() Logger {
	return l
}

type StdoutLogger struct{}

func (l *StdoutLogger) Printf(format string, v ...interface{}) {
	output := fmt.Sprintf(format, v...)
	fmt.Println(output)
}
