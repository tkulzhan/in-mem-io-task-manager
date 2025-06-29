package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	internal zerolog.Logger
}

func New() *Logger {
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
