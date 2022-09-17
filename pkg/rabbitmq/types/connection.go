package types

import (
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/rabbitmq/amqp091-go"
)

type connection struct {
	cfg *config.RabbitMQConfig
	*amqp091.Connection
}

type IConnection interface {
	IsClosed() bool
	// Channel gets a new channel on this connection
	Channel() (*amqp091.Channel, error)
	Close() error
	NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error
	Reconnect() (IConnection, error)
	Raw() *amqp091.Connection
}

func NewConnection(cfg *config.RabbitMQConfig) (IConnection, error) {
	if cfg.RabbitMqHostOptions == nil {
		return nil, errors.New("RabbitMqHostOptions can't be nil in RabbitMQConfig")
	}

	conn, err := amqp091.Dial(cfg.RabbitMqHostOptions.EndPoint())
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &connection{
		cfg:        cfg,
		Connection: conn,
	}, nil
}

func (c *connection) Reconnect() (IConnection, error) {
	return NewConnection(c.cfg)
}

func (c *connection) Raw() *amqp091.Connection {
	return c.Connection
}

func (c *connection) Channel() (*amqp091.Channel, error) {
	return c.Connection.Channel()
}
