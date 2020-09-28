package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	Logger = logger
}
