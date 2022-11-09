package rabbitmq

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/contracts"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

type rabbitmqTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.RabbitMQContainerOptions
}

func NewRabbitMQTestContainers() *rabbitmqTestContainers {
	return &rabbitmqTestContainers{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672/tcp", "15672/tcp", "15671/tcp", "25672/tcp", "5671/tcp"},
			Host:        "localhost",
			VirtualHost: "/",
			UserName:    "guest",
			Password:    "guest",
			Tag:         "3-management",
			ImageName:   "rabbitmq",
			Name:        "rabbitmq-testcontainers",
		},
	}
}

func (g *rabbitmqTestContainers) Start(ctx context.Context, t *testing.T, rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc, options ...*contracts.RabbitMQContainerOptions) (bus.Bus, error) {
	//https://github.com/testcontainers/testcontainers-go
	//https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

	//TODO: Using Parallel Container
	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// get a free random host hostPort
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[0]))
	if err != nil {
		return nil, err
	}

	uiHttpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[1]))
	if err != nil {
		return nil, err
	}
	t.Logf("rabbitmq ui port is: %d", uiHttpPort.Int())

	g.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		_ = dbContainer.Terminate(ctx)
	})

	mqBus, err := bus2.NewRabbitMQBus(
		ctx,
		&config.RabbitMQConfig{
			RabbitMqHostOptions: &config.RabbitMqHostOptions{
				UserName:    g.defaultOptions.UserName,
				Password:    g.defaultOptions.Password,
				HostName:    host,
				VirtualHost: g.defaultOptions.VirtualHost,
				Port:        g.defaultOptions.HostPort,
			},
		},
		rabbitmqBuilderFunc,
		json.NewJsonEventSerializer(),
		defaultLogger.Logger)
	if err != nil {
		return nil, err
	}

	return mqBus, nil
}

func (g *rabbitmqTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
}

func (g *rabbitmqTestContainers) getRunOptions(opts ...*contracts.RabbitMQContainerOptions) testcontainers.ContainerRequest {
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
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS": g.defaultOptions.Password,
		},
	}

	return containerReq
}
