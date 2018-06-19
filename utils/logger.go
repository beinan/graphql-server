package utils

import (
	"go.uber.org/zap"
	"log"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
}

var DefaultLogger = NewLogger()

func NewLogger() Logger {
	var logger *zap.Logger
	var err error
	if IsProd() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync() // flushes buffer
	sugar := logger.Sugar()
	return sugar
}
