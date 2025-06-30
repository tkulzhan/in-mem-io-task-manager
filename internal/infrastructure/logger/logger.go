package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger struct {
	internal zerolog.Logger
}

func New(logLevel string) *Logger {
	zerolog.SetGlobalLevel(parseLogLevel(logLevel))
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{
		internal: zlog,
	}
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.internal.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	event := l.internal.Error()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.internal.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (l *Logger) Trace(msg string, fields map[string]interface{}) {
	event := l.internal.Trace()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
