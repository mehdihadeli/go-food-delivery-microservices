package fxapp

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/fx"

	"github.com/spf13/viper"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	loggerConfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

type applicationBuilder struct {
	provides    []interface{}
	decorates   []interface{}
	options     []fx.Option
	logger      logger.Logger
	environment environemnt.Environment
}

func NewApplicationBuilder(environments ...environemnt.Environment) contracts.ApplicationBuilder {
	env := environemnt.ConfigAppEnv(environments...)

	setConfigPath()

	var logger logger.Logger
	logoption, err := loggerConfig.ProvideLogConfig(env)
	if err != nil || logoption == nil {
		logger = zap.NewZapLogger(logoption, env)
	} else if logoption.LogType == models.Logrus {
		logger = logrous.NewLogrusLogger(logoption, env)
	} else {
		logger = zap.NewZapLogger(logoption, env)
	}

	return &applicationBuilder{logger: logger, environment: env}
}

func setConfigPath() {
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// get from `os env` or viper `internal registry`
	pn := viper.Get(constants.PROJECT_NAME_ENV)
	if pn == nil {
		return
	}
	for !strings.HasSuffix(wd, pn.(string)) {
		wd = filepath.Dir(wd)
	}

	// Get the absolute path of the executed project directory
	absCurrentDir, _ := filepath.Abs(wd)

	viper.Set(constants.AppRootPath, absCurrentDir)

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")

	viper.Set(constants.ConfigPath, configPath)
}

func (a *applicationBuilder) ProvideModule(module fx.Option) {
	a.options = append(a.options, module)
}

func (a *applicationBuilder) Provide(constructors ...interface{}) {
	a.provides = append(a.provides, constructors...)
}

func (a *applicationBuilder) Decorate(constructors ...interface{}) {
	a.decorates = append(a.decorates, constructors...)
}

func (a *applicationBuilder) Build() contracts.Application {
	app := NewApplication(a.provides, a.decorates, a.options, a.logger, a.environment)

	return app
}

func (a *applicationBuilder) GetProvides() []interface{} {
	return a.provides
}

func (a *applicationBuilder) GetDecorates() []interface{} {
	return a.decorates
}

func (a *applicationBuilder) Options() []fx.Option {
	return a.options
}

func (a *applicationBuilder) Logger() logger.Logger {
	return a.logger
}

func (a *applicationBuilder) Environment() environemnt.Environment {
	return a.environment
}
