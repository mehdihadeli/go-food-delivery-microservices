package config

import (
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
)

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"ServiceName"`
}

func NewAppConfig(env config.Environment) (*AppOptions, error) {
	cfg, err := config.BindConfig[*AppOptions](env.GetEnvironmentName())
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
