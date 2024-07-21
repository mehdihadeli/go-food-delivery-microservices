package pipelines

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
)

type config struct {
	logger      logger.Logger
	serviceName string
}

var defaultConfig = &config{
	serviceName: "app",
	logger:      defaultLogger.GetLogger(),
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}

func WithLogger(l logger.Logger) Option {
	return optionFunc(func(cfg *config) {
		if cfg.logger != nil {
			cfg.logger = l
		}
	})
}
