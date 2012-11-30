package protocol

type Logger interface {
	Printf(format string, v ...interface{})
}

type NullLogger struct{}

func (l *NullLogger) Printf(format string, v ...interface{}) {}
