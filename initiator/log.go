package initiator

import (
	"go.uber.org/zap"
)

func InitialzeLog() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger
}
