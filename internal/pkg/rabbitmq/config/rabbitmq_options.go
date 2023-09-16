package config

import (
	"fmt"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/iancoleman/strcase"
)

type RabbitmqOptions struct {
	RabbitmqHostOptions *RabbitmqHostOptions `mapstructure:"rabbitmqHostOptions"`
	DeliveryMode        uint8
	Persisted           bool
	AppId               string
	AutoStart           bool `mapstructure:"autoStart"           default:"true"`
	Reconnecting        bool `mapstructure:"reconnecting"        default:"true"`
}

type RabbitmqHostOptions struct {
	HostName    string    `mapstructure:"hostName"`
	VirtualHost string    `mapstructure:"virtualHost"`
	Port        int       `mapstructure:"port"`
	HttpPort    int       `mapstructure:"httpPort"`
	UserName    string    `mapstructure:"userName"`
	Password    string    `mapstructure:"password"`
	RetryDelay  time.Time `mapstructure:"retryDelay"`
}

func (h *RabbitmqHostOptions) AmqpEndPoint() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", h.UserName, h.Password, h.HostName, h.Port)
}

func (h *RabbitmqHostOptions) HttpEndPoint() string {
	return fmt.Sprintf("http://%s:%d", h.HostName, h.HttpPort)
}

func ProvideConfig(environment environemnt.Environment) (*RabbitmqOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetTypeNameByT[RabbitmqOptions]())
	cfg, err := config.BindConfigKey[*RabbitmqOptions](optionName, environment)

	return cfg, err
}
