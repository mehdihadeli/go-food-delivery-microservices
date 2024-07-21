package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/contracts"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// https://github.com/testcontainers/testcontainers-go/issues/1359
// https://github.com/testcontainers/testcontainers-go/issues/1249

type rabbitmqTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.RabbitMQContainerOptions
	logger         logger.Logger
}

func NewRabbitMQTestContainers(l logger.Logger) contracts.RabbitMQContainer {
	return &rabbitmqTestContainers{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672/tcp", "15672/tcp"},
			Host:        "localhost",
			VirtualHost: "/",
			UserName:    "guest",
			Password:    "guest",
			HttpPort:    15672,
			HostPort:    5672,
			Tag:         "management",
			ImageName:   "rabbitmq",
			Name:        "rabbitmq-testcontainers",
		},
		logger: l,
	}
}

func (g *rabbitmqTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqHostOptions, error) {
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if terr := dbContainer.Terminate(ctx); terr != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		time.Sleep(time.Second * 1)
	})

	// get a free random host port for rabbitmq `Tcp Port`
	hostPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port(g.defaultOptions.Ports[0]),
	)
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()
	g.logger.Infof("rabbitmq host port is: %d", hostPort.Int())

	// https://github.com/michaelklishin/rabbit-hole/issues/74
	// get a free random host port for rabbitmq UI `Http Port`
	uiHttpPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port(g.defaultOptions.Ports[1]),
	)
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HttpPort = uiHttpPort.Int()
	g.logger.Infof("rabbitmq ui port is: %d", uiHttpPort.Int())

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	isConnectable := IsConnectable(g.logger, g.defaultOptions)
	if !isConnectable {
		return g.PopulateContainerOptions(context.Background(), t, options...)
	}

	g.container = dbContainer

	option := &config.RabbitmqHostOptions{
		UserName:    g.defaultOptions.UserName,
		Password:    g.defaultOptions.Password,
		HostName:    host,
		VirtualHost: g.defaultOptions.VirtualHost,
		Port:        g.defaultOptions.HostPort,
		HttpPort:    g.defaultOptions.HttpPort,
	}

	return option, nil
}

func (g *rabbitmqTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

func (g *rabbitmqTestContainers) getRunOptions(
	opts ...*contracts.RabbitMQContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if len(option.Ports) > 0 {
			g.defaultOptions.Ports = option.Ports
		}
		if option.UserName != "" {
			g.defaultOptions.UserName = option.UserName
		}
		if option.Password != "" {
			g.defaultOptions.Password = option.Password
		}
		if option.Tag != "" {
			g.defaultOptions.Tag = option.Tag
		}
	}
	containerReq := testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			g.defaultOptions.ImageName,
			g.defaultOptions.Tag,
		),
		ExposedPorts: g.defaultOptions.Ports,
		WaitingFor: wait.ForListeningPort(
			nat.Port(g.defaultOptions.Ports[0]),
		),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Hostname: g.defaultOptions.Host,
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS": g.defaultOptions.Password,
		},
	}

	return containerReq
}

func IsConnectable(
	logger logger.Logger,
	options *contracts.RabbitMQContainerOptions,
) bool {
	conn, err := amqp091.Dial(options.AmqpEndPoint())
	if err != nil {
		logError(
			logger,
			options.UserName,
			options.Password,
			options.Host,
			options.HostPort,
		)

		return false
	}

	defer conn.Close()

	if err != nil || (conn != nil && conn.IsClosed()) {
		logError(
			logger,
			options.UserName,
			options.Password,
			options.Host,
			options.HostPort,
		)

		return false
	}

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		options.HttpEndPoint(),
		options.UserName,
		options.Password,
	)
	_, err = rmqc.ListExchanges()

	if err != nil {
		logger.Errorf(
			"Error in creating rabbitmq connection with http host: %s",
			options.HttpEndPoint(),
		)

		return false
	}

	logger.Infof(
		"Opened rabbitmq connection on host: amqp://%s:%s@%s:%d",
		options.UserName,
		options.Password,
		options.Host,
		options.HostPort,
	)

	return true
}

func logError(
	logger logger.Logger,
	userName string,
	password string,
	host string,
	hostPort int,
) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf(
		"Error in creating rabbitmq connection with amqp host: amqp://%s:%s@%s:%d",
		userName,
		password,
		host,
		hostPort,
	)
}
