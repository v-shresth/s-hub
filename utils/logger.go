package utils

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() Logger {
	return Logger{
		log: logrus.New(),
	}
}

func (l *Logger) WithError(err error) *Logger {
	l.log.WithError(err)
	return l
}

func (l *Logger) Error(message string) *Logger {
	l.log.Error(message)
	return l
}

func (l *Logger) Info(message string) *Logger {
	l.log.Info(message)
	return l
}

func (l *Logger) Debug(message string) *Logger {
	l.log.Debug(message)
	return l
}

func (l *Logger) Warn(message string) *Logger {
	l.log.Warn(message)
	return l
}

func (l *Logger) Fatal(message string) *Logger {
	l.log.Fatal(message)
	return l
}

func (l *Logger) Panic(message string) *Logger {
	l.log.Panic(message)
	return l
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.log.WithField(key, value)
	return l
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.log.WithFields(fields)
	return l
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	l.log.WithContext(ctx)
	return l
}
