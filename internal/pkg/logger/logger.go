package logger

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	return log
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf logs an error message with format
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof logs an info message with format
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
