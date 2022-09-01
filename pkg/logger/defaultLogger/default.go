package defaultLogger

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/logrous"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/zap"
	"os"
)

var (
	Logger logger.Logger
)

func init() {
	logType := os.Getenv("LogConfig_LogType")
	if logType == "" {
		logType = "Zap"
	}
	switch logType {
	case "Zap":
		Logger = zap.NewZapLogger(&logger.LogConfig{
			LogLevel: "debug",
			LogType:  logger.Zap,
		})
		break
	case "Logrus":
		Logger = logrous.NewLogrusLogger(&logger.LogConfig{
			LogLevel: "debug",
			LogType:  logger.Logrus,
		})
		break
	}
}
