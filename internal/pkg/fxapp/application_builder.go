package fxapp

import (
	"go.uber.org/fx"

	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

type ApplicationBuilder struct {
	Providers []interface{}
	Options   []fx.Option
	Logger    logger.Logger
}

func NewApplicationBuilder() *ApplicationBuilder {
	var logger logger.Logger
	env := config2.ConfigAppEnv()
	logoption, err := config.ProvideLogConfig()
	if err != nil || logoption == nil {
		logger = zap.NewZapLogger(logoption, env)
	} else if logoption.LogType == models.Logrus {
		logger = logrous.NewLogrusLogger(logoption, env)
	} else {
		logger = zap.NewZapLogger(logoption, env)
	}

	return &ApplicationBuilder{Logger: logger}
}

func (a *ApplicationBuilder) ProvideModule(module fx.Option) {
	a.Options = append(a.Options, module)
}

func (a *ApplicationBuilder) Provide(constructors ...interface{}) {
	a.Providers = append(a.Providers, constructors...)
}

func (a *ApplicationBuilder) Build() *Application {
	app := NewApplication(a.Providers, a.Options, a.Logger)

	return app
}
