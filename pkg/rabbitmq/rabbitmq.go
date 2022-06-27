package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Exchange       string
	Queue          string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int

}

func NewRabbitMQConn(cfg *RabbitMQConfig) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	return amqp.Dial(connAddr)
}
