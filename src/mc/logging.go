package mc

import (
	"fmt"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type StdoutLogger struct{}

func (l *StdoutLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
