package ax

import (
	"fmt"
	"runtime"
	"strings"
)

type traceableError struct {
	message    string
	stacktrace []string
}

func (e *traceableError) Error() string {
	tabulated_stacktrace := strings.Join(e.stacktrace, "\n\t")
	return fmt.Sprintf("%s:\n\t%s", e.message, tabulated_stacktrace)
}

func getStackTrace(offset int) []string {
	strstack := make([]string, 0)
	stack := make([]uintptr, 100)
	count := runtime.Callers(offset, stack)
	stack = stack[0:count]
	for i, pc := range stack {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		if i == 0 {
			line--
		}
		strstack = append(strstack, fmt.Sprintf("%s:%d\n    inside %s", file, line, fn.Name()))
	}
	return strstack
}

func coerceToError(msg string, stackOffset int) error {
	return &traceableError{
		message:    msg,
		stacktrace: getStackTrace(stackOffset),
	}
}

func WrapError(err error) error {
	if err != nil {
		return coerceToError(err.Error(), 3)
	}
	return nil
}

func Error(message string) error {
	return coerceToError(message, 3)
}

func Errorf(format string, v ...interface{}) error {
	return coerceToError(fmt.Sprintf(format, v...), 4)
}
