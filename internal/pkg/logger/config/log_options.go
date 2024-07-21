package config

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/models"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[LogOptions]())

type LogOptions struct {
	LogLevel      string         `mapstructure:"level"`
	LogType       models.LogType `mapstructure:"logType"`
	CallerEnabled bool           `mapstructure:"callerEnabled"`
	EnableTracing bool           `mapstructure:"enableTracing" default:"true"`
}

func ProvideLogConfig(env environment.Environment) (*LogOptions, error) {
	return config.BindConfigKey[*LogOptions](optionName, env)
}
