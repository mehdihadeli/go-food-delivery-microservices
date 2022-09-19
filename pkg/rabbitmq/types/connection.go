package types

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/rabbitmq/amqp091-go"
)

type internalConnection struct {
	cfg *config.RabbitMQConfig
	*amqp091.Connection
	isConnected     bool
	errChan         chan error
	reconnectedChan chan struct{}
}

type IConnection interface {
	IsClosed() bool
	IsConnected() bool
	// Channel gets a new channel on this internalConnection
	Channel() (*amqp091.Channel, error)
	Close() error
	ReConnect() error
	NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error
	Raw() *amqp091.Connection
	ErrorChannel() chan error
	ReconnectedChannel() chan struct{}
}

func NewConnection(ctx context.Context, cfg *config.RabbitMQConfig) (IConnection, error) {
	//https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	if cfg.RabbitMqHostOptions == nil {
		return nil, errors.New("rabbitmq host options is nil")
	}

	c := &internalConnection{
		cfg:             cfg,
		errChan:         make(chan error),
		reconnectedChan: make(chan struct{}),
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	go c.handleReconnecting(ctx)

	return c, err
}

func (c *internalConnection) Close() error {
	return c.Connection.Close()
}

func (c *internalConnection) IsConnected() bool {
	return c.isConnected
}

func (c *internalConnection) ErrorChannel() chan error {
	return c.errChan
}

func (c *internalConnection) ReconnectedChannel() chan struct{} {
	return c.reconnectedChan
}

func (c *internalConnection) ReConnect() error {
	if c.Connection.IsClosed() == false {
		return nil
	}

	return c.connect()
}

func (c *internalConnection) Raw() *amqp091.Connection {
	return c.Connection
}

func (c *internalConnection) Channel() (*amqp091.Channel, error) {
	return c.Connection.Channel()
}

func (c *internalConnection) connect() error {
	conn, err := amqp091.Dial(c.cfg.RabbitMqHostOptions.EndPoint())
	if err != nil {
		return errors.WrapIf(err, fmt.Sprintf("Error in creating rabbitmq connection with %s", c.cfg.RabbitMqHostOptions.EndPoint()))
	}

	c.Connection = conn
	c.isConnected = true

	//https://stackoverflow.com/questions/41991926/how-to-detect-dead-rabbitmq-connection
	notifyClose := c.Connection.NotifyClose(make(chan *amqp091.Error))

	go func() {
		<-notifyClose //Listen to NotifyClose
		c.isConnected = false
		c.errChan <- errors.New("Connection Closed")
	}()

	return nil
}

func (c *internalConnection) handleReconnecting(ctx context.Context) {
	for {
		select {
		case err := <-c.errChan:
			if err != nil {
				defaultLogger.Logger.Info("Rabbitmq Connection Reconnecting started")
				err := c.connect()
				if err != nil {
					continue
				}
				defaultLogger.Logger.Info("Rabbitmq Connection Reconnected")
				c.isConnected = true
				c.reconnectedChan <- struct{}{}
				continue
			}
		case <-ctx.Done():
			c.Close()
			return
		}
	}
}
