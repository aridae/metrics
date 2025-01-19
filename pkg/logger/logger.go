package logger

import (
	"log"

	"go.uber.org/zap"
)

var (
	_logger *zap.SugaredLogger
)

func init() {
	devLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	_logger = devLogger.Sugar()
}

func Fatalf(template string, args ...interface{}) {
	_logger.Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_logger.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	_logger.Warnf(template, args...)
}

func Infof(template string, args ...interface{}) {
	_logger.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	_logger.Debugf(template, args...)
}
