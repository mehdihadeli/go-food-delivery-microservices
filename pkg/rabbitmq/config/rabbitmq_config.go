package config

import (
	"fmt"
	"time"
)

type RabbitMQConfig struct {
	RabbitMqHostOptions *RabbitMqHostOptions
	DeliveryMode        uint8
	Persisted           bool
	AppId               string
}

type RabbitMqHostOptions struct {
	HostName    string    `mapstructure:"hostName"`
	VirtualHost string    `mapstructure:"virtualHost"`
	Port        int       `mapstructure:"port"`
	UserName    string    `mapstructure:"userName"`
	Password    string    `mapstructure:"password"`
	RetryDelay  time.Time `mapstructure:"retryDelay"`
}

func (h *RabbitMqHostOptions) EndPoint() string {
	return fmt.Sprintf("amqp://%s:%d", h.HostName, h.Port)
}
