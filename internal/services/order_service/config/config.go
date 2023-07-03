package config

import (
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
)

type Config struct {
	AppOptions AppOptions `mapstructure:"appOptions"`
}

func NewConfig(environment environemnt.Environment) (*Config, error) {
	cfg, err := config.BindConfig[*Config](environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type AppOptions struct {
	DeliveryType string `mapstructure:"deliveryType"`
	ServiceName  string `mapstructure:"serviceName"`
}

func (cfg *AppOptions) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppOptions) GetMicroserviceName() string {
	return cfg.ServiceName
}
