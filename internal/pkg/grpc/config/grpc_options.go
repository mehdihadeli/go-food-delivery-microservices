package config

import (
	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[GrpcOptions]())

type GrpcOptions struct {
	Port        string `mapstructure:"port"        env:"Port"`
	Host        string `mapstructure:"host"        env:"Host"`
	Development bool   `mapstructure:"development" env:"Development"`
	Name        string `mapstructure:"name"        env:"Name"`
}

func ProvideConfig() (*GrpcOptions, error) {
	return config.BindConfigKey[*GrpcOptions](optionName)
}
