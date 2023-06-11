//go:build.sh go1.18

package config

import (
	"fmt"
	"time"

	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[RabbitmqOptions]())

type RabbitmqOptions struct {
	RabbitmqHostOptions *RabbitmqHostOptions `mapstructure:"rabbitmqHostOptions"`
	DeliveryMode        uint8
	Persisted           bool
	AppId               string
}

type RabbitmqHostOptions struct {
	HostName    string    `mapstructure:"hostName"`
	VirtualHost string    `mapstructure:"virtualHost"`
	Port        int       `mapstructure:"port"`
	UserName    string    `mapstructure:"userName"`
	Password    string    `mapstructure:"password"`
	RetryDelay  time.Time `mapstructure:"retryDelay"`
}

func (h *RabbitmqHostOptions) EndPoint() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", h.UserName, h.Password, h.HostName, h.Port)
}

func ProvideConfig() (*RabbitmqOptions, error) {
	return config.BindConfigKey[*RabbitmqOptions](optionName)
}
