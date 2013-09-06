package ax

// A NullLogger instance. A quick way to access the NullLogger
var nullLogger Logger = &NullLogger{}

// The core interface for all logging
type Logger interface {
	Printf(format string, values ...interface{})
}

// The interface for a logger that can wrap another logger
type WrapLogger interface {
	Logger
	SetLogger(l Logger)
	WrappedLogger() Logger
}

// Shorthand for ensuring a logger is returned
// You can provide as many fallback loggers as you want.
// If all loggers are nil, the ax.nullLogger logger is used.
func Use(loggers ...Logger) Logger {
	for _, logger := range loggers {
		if logger != nil {
			return logger
		}
	}
	return nullLogger
}

// Shorthand for wrapping a logger with WrapLoggers.
// Multiple WrapLoggers can be provided to create a nested behavior
// with the inner-most WrapLogger placed before the all other WrapLoggers.
func Wrap(logger Logger, outerLoggers ...WrapLogger) WrapLogger {
	previousLogger := logger
	for _, outerLogger := range outerLoggers {
		outerLogger.SetLogger(previousLogger)
		previousLogger = outerLogger
	}
	return previousLogger.(WrapLogger)
}
