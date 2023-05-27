package defaultLogger

import (
    "os"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

var (
    Logger logger.Logger
)

func init() {
    logType := os.Getenv("LogConfig_LogType")

    switch logType {
    case "Zap", "":
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
    default:

    }
}
