package defaultLogger

import (
	"os"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

var Logger logger.Logger

func SetupDefaultLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		Logger = zap.NewZapLogger(
			&config.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.Dev,
		)
		break
	case "Logrus":
		Logger = logrous.NewLogrusLogger(
			&config.LogOptions{LogType: models.Logrus, CallerEnabled: false},
			constants.Dev,
		)
		break
	default:
	}
}
