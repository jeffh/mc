package ax

var DefaultLogger Logger = &NullLogger{}

type Logger interface {
	Printf(format string, values ...interface{})
}

type WrapLogger interface {
	Logger
	SetLogger(l Logger)
	WrappedLogger() Logger
}

func UseOrDefault(logger Logger, defaultLogger Logger) Logger {
	if logger == nil {
		return defaultLogger
	}
	return logger
}

func Use(logger Logger) Logger {
	return UseOrDefault(logger, DefaultLogger)
}

func Wrap(logger Logger, outerLoggers ...WrapLogger) WrapLogger {
	previousLogger := logger
	for _, outerLogger := range outerLoggers {
		outerLogger.SetLogger(previousLogger)
		previousLogger = outerLogger
	}
	return previousLogger.(WrapLogger)
}
