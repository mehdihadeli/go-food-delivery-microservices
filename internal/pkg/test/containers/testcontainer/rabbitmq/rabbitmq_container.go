package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"
)

type rabbitmqTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.RabbitMQContainerOptions
}

func NewRabbitMQTestContainers() contracts.RabbitMQContainer {
	return &rabbitmqTestContainers{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672/tcp", "15672/tcp", "15671/tcp", "25672/tcp", "5671/tcp"},
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
	}
}

func (g *rabbitmqTestContainers) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqOptions, error) {
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

	// TODO: Using Parallel Container
	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// get a free random host port for rabbitmq `Tcp Port`
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[0]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()

	// https://github.com/michaelklishin/rabbit-hole/issues/74
	// get a free random host port for rabbitmq UI `Http Port`
	uiHttpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[1]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HttpPort = uiHttpPort.Int()
	t.Logf("rabbitmq ui port is: %d", uiHttpPort.Int())

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		_ = dbContainer.Terminate(ctx)
	})

	option := &config.RabbitmqOptions{
		RabbitmqHostOptions: &config.RabbitmqHostOptions{
			UserName:    g.defaultOptions.UserName,
			Password:    g.defaultOptions.Password,
			HostName:    host,
			VirtualHost: g.defaultOptions.VirtualHost,
			Port:        g.defaultOptions.HostPort,
			HttpPort:    g.defaultOptions.HttpPort,
		},
	}

	return option, nil
}

func (g *rabbitmqTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	serializer serializer.EventSerializer,
	logger logger.Logger,
	rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
	options ...*contracts.RabbitMQContainerOptions,
) (bus.Bus, error) {
	rabbitOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	mqBus, err := bus2.NewRabbitmqBus(
		rabbitOptions,
		serializer,
		logger,
		rabbitmqBuilderFunc,
	)
	if err != nil {
		return nil, err
	}

	return mqBus, nil
}

func (g *rabbitmqTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
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
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: g.defaultOptions.Ports,
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Ports[0])),
		Hostname:     g.defaultOptions.Host,
		SkipReaper:   true,
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS": g.defaultOptions.Password,
		},
	}

	return containerReq
}
