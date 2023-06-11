package defaultLogger

import (
	"log"
	"os"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

var Logger logger.Logger

func init() {
	logType := os.Getenv("LogConfig_LogType")

	cfg, err := config.ProvideLogConfig()
	if err != nil {
		log.Fatal(err)
	}

	switch logType {
	case "Zap", "":
		Logger = zap.NewZapLogger(cfg, constants.Dev)
		break
	case "Logrus":
		Logger = logrous.NewLogrusLogger(cfg, constants.Dev)
		break
	default:
	}
}
