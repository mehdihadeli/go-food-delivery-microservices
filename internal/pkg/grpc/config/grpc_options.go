package config

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GrpcOptions]())

type GrpcOptions struct {
	Port        string `mapstructure:"port"        env:"TcpPort"`
	Host        string `mapstructure:"host"        env:"Host"`
	Development bool   `mapstructure:"development" env:"Development"`
	Name        string `mapstructure:"name"        env:"ShortTypeName"`
}

func ProvideConfig(environment environment.Environment) (*GrpcOptions, error) {
	return config.BindConfigKey[*GrpcOptions](optionName, environment)
}
