package config

import (
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
)

type Config struct {
	AppOptions AppOptions `mapstructure:"appOptions" env:"AppOptions"`
}

func NewConfig(env environemnt.Environment) (*Config, error) {
	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType" env:"DeliveryType"`
	ServiceName  string `mapstructure:"serviceName"  env:"ServiceName"`
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
