package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger   *zap.Logger
	syncOnce sync.Once
)

func InitLogger() error {
	var err error
	syncOnce.Do(func() {
		logger, err = zap.NewProduction()
	})
	return err
}

func GetLogger() *zap.Logger {
	return logger
}
