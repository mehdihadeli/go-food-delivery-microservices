package types

import (
	"fmt"

	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	errorUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils/errorutils"

	"emperror.dev/errors"
	"github.com/rabbitmq/amqp091-go"
)

type internalConnection struct {
	cfg *config.RabbitmqOptions
	*amqp091.Connection
	isConnected       bool
	errConnectionChan chan error
	errChannelChan    chan error
	reconnectedChan   chan struct{}
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
	ErrorConnectionChannel() chan error
	ReconnectedChannel() chan struct{}
}

func NewRabbitMQConnection(cfg *config.RabbitmqOptions) (IConnection, error) {
	// https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	if cfg.RabbitmqHostOptions == nil {
		return nil, errors.New("rabbitmq host options is nil")
	}

	c := &internalConnection{
		cfg:               cfg,
		errConnectionChan: make(chan error),
		// errChannelChan:    make(chan error),
		reconnectedChan: make(chan struct{}),
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	if cfg.Reconnecting {
		go c.handleReconnecting()
	}

	return c, err
}

func (c *internalConnection) Close() error {
	return c.Connection.Close()
}

func (c *internalConnection) IsConnected() bool {
	return c.isConnected
}

func (c *internalConnection) ErrorConnectionChannel() chan error {
	return c.errConnectionChan
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
	ch, err := c.Connection.Channel()
	//notifyChannelClose := ch.NotifyClose(make(chan *amqp091.Error))
	//go func() {
	//	<-notifyChannelClose //Listen to notifyChannelClose
	//	c.errChannelChan <- errors.New("Channel Closed")
	//}()

	return ch, err
}

func (c *internalConnection) connect() error {
	conn, err := amqp091.Dial(c.cfg.RabbitmqHostOptions.AmqpEndPoint())
	if err != nil {
		return errors.WrapIf(
			err,
			fmt.Sprintf(
				"Error in connecting to rabbitmq with host: %s",
				c.cfg.RabbitmqHostOptions.AmqpEndPoint(),
			),
		)
	}

	c.Connection = conn
	c.isConnected = true

	// https://stackoverflow.com/questions/41991926/how-to-detect-dead-rabbitmq-connection
	notifyClose := c.Connection.NotifyClose(make(chan *amqp091.Error))

	go func() {
		defer errorUtils.HandlePanic()
		chanErr := <-notifyClose // Listen to NotifyClose
		c.isConnected = false
		c.errConnectionChan <- errors.WrapIf(chanErr, "rabbitmq Connection Closed with an error.")
	}()

	return nil
}

func (c *internalConnection) handleReconnecting() {
	defer errorUtils.HandlePanic()
	for {
		select {
		case err := <-c.errConnectionChan:
			if err != nil {
				defaultLogger.GetLogger().
					Info("Rabbitmq Connection Reconnecting started")
				err := c.connect()
				if err != nil {
					defaultLogger.GetLogger().
						Error(fmt.Sprintf("Error in reconnecting, %s", err))
					continue
				}

				defaultLogger.GetLogger().
					Info("Rabbitmq Connection Reconnected")
				c.isConnected = true
				c.reconnectedChan <- struct{}{}
				continue
			}
		}
	}
}
