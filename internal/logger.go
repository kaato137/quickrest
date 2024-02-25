package internal

import (
	"os"

	"github.com/charmbracelet/log"
)

type Logger interface {
	Debug(msg any, keyvals ...any)
	Debugf(format string, keyvals ...any)
	Info(msg any, keyvals ...any)
	Infof(format string, keyvals ...any)
	Warn(msg any, keyvals ...any)
	Warnf(format string, keyvals ...any)
	Error(msg any, keyvals ...any)
	Errorf(format string, keyvals ...any)
	Fatal(msg any, keyvals ...any)
	Fatalf(format string, keyvals ...any)
}

type CustomLogger struct {
	logger *log.Logger
}

func NewLogger() *CustomLogger {
	return &CustomLogger{
		logger: log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: true,
		}),
	}
}

func (l *CustomLogger) Debug(msg any, keyvals ...any) {
	l.logger.Debug(msg, keyvals...)
}

func (l *CustomLogger) Debugf(format string, keyvals ...any) {
	l.logger.Debugf(format, keyvals...)
}

func (l *CustomLogger) Info(msg any, keyvals ...any) {
	l.logger.Info(msg, keyvals...)
}

func (l *CustomLogger) Infof(format string, keyvals ...any) {
	l.logger.Infof(format, keyvals...)
}

func (l *CustomLogger) Warn(msg any, keyvals ...any) {
	l.logger.Warn(msg, keyvals...)
}

func (l *CustomLogger) Warnf(format string, keyvals ...any) {
	l.logger.Warnf(format, keyvals...)
}

func (l *CustomLogger) Error(msg any, keyvals ...any) {
	l.logger.Error(msg, keyvals...)
}

func (l *CustomLogger) Errorf(format string, keyvals ...any) {
	l.logger.Errorf(format, keyvals...)
}

func (l *CustomLogger) Fatal(msg any, keyvals ...any) {
	l.logger.Fatal(msg, keyvals...)
}

func (l *CustomLogger) Fatalf(format string, keyvals ...any) {
	l.logger.Fatalf(format, keyvals...)
}
