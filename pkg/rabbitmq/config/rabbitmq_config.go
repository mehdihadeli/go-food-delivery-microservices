package config

import (
	"fmt"
	"time"
)

type RabbitMQConfig struct {
	RabbitMqHostOptions *RabbitMqHostOptions
	DeliveryMode        bool
	Persisted           bool
	AppId               string
}

type RabbitMqHostOptions struct {
	HostName    string
	VirtualHost string
	Port        int
	UserName    string
	Password    string
	RetryDelay  time.Time
}

func (h *RabbitMqHostOptions) EndPoint() string {
	return fmt.Sprintf("amqp://%s:%d", h.HostName, h.Port)
}
