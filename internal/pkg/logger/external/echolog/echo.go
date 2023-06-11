package echolog

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Ref: https://articles.wesionary.team/logging-interfaces-in-go-182c28be3d18

type EchoCustomLogger struct {
	logger.Logger
}

func NewEchoCustomLogger(logger logger.Logger) *EchoCustomLogger {
	//cfg, err := config.ProvideLogConfig()
	//
	//var logger logger.Logger
	//if cfg.LogType == models.Logrus && err != nil {
	//	logger = logrous.NewLogrusLogger(cfg, constants.Dev)
	//} else {
	//	if err != nil {
	//		cfg = &config.LogOptions{LogLevel: "info", LogType: models.Zap}
	//	}
	//	logger = zap.NewZapLogger(cfg, constants.Dev)
	//}

	return &EchoCustomLogger{
		Logger: logger,
	}
}
