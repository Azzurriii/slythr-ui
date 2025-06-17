package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

type Logger struct{}

var Default = &Logger{}

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

func Infof(format string, args ...interface{}) {
	log.Info(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	log.Error(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	log.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	Infof(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	Warnf(format, args...)
}
