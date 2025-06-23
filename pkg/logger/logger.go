package logger

import (
	"os"
	"path/filepath"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// Logger wraps gookit/slog for better compatibility
type Logger struct {
	*slog.Logger
}

var Default *Logger

func init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		consoleHandler := handler.NewConsoleHandler(slog.AllLevels)
		logger := slog.NewWithHandlers(consoleHandler)
		Default = &Logger{Logger: logger}
		return
	}

	logFile := filepath.Join("logs", "app.log")
	fileHandler, err := handler.NewFileHandler(logFile, handler.WithLogLevels(slog.AllLevels))
	if err != nil {
		consoleHandler := handler.NewConsoleHandler(slog.AllLevels)
		logger := slog.NewWithHandlers(consoleHandler)
		Default = &Logger{Logger: logger}
		return
	}

	consoleHandler := handler.NewConsoleHandler(slog.AllLevels)

	logger := slog.NewWithHandlers(fileHandler, consoleHandler)

	Default = &Logger{Logger: logger}
}

func NewLogger() *Logger {
	return Default
}

func NewLoggerWithLevel(level slog.Level) *Logger {
	h := handler.NewConsoleHandler([]slog.Level{level, slog.ErrorLevel, slog.FatalLevel})
	logger := slog.NewWithHandlers(h)
	return &Logger{Logger: logger}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	Default.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

func Error(args ...interface{}) {
	Default.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	Default.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	Default.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}
