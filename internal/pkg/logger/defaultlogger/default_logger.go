package defaultLogger

import (
	"os"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/models"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
)

var l logger.Logger

func initLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		l = zap.NewZapLogger(
			&config.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.Dev,
		)
		break
	case "Logrus":
		l = logrous.NewLogrusLogger(
			&config.LogOptions{LogType: models.Logrus, CallerEnabled: false},
			constants.Dev,
		)
		break
	default:
	}
}

func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
