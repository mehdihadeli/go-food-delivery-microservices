package config

import (
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[AppOptions]())

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"ServiceName"`
}

func NewAppConfig(env config.Environment) (*AppOptions, error) {
	cfg, err := config.BindConfigKey[*AppOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
