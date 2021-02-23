package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger represent common interface for logging function
type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

var log Logger

func init() {
	// default is logrus
	log = &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

// SetLogger will set an logger implements for the sdk
func SetLogger(l Logger) {
	log = l
}

// Errorf log error
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Error log error
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatalf log fatal
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Fatal log fatal
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Infof log info
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Info log info
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warnf log warn
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Warn log warn
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Debugf log debug
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Debug log debug
func Debug(args ...interface{}) {
	log.Debug(args...)
}
