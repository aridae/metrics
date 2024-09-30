package logger

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var (
	_logger *zap.SugaredLogger
	_once   sync.Once
)

func Obtain() *zap.SugaredLogger {
	_once.Do(func() {
		devLogger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}

		_logger = devLogger.Sugar()
	})

	return _logger
}
