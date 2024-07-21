package testcontainer

import (
	"context"
	"fmt"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
) (*config.RabbitmqHostOptions, error) {
	t.Helper()

	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := testcontainers.ContainerRequest{
		Image: fmt.Sprintf(
			"%s:%s",
			"rabbitmq",
			"management",
		),
		ExposedPorts: []string{"5672/tcp", "15672/tcp"},
		WaitingFor: wait.ForListeningPort(
			nat.Port("5672/tcp"),
		),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Hostname: "localhost",
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": "guest",
			"RABBITMQ_DEFAULT_PASS": "guest",
		},
	}

	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	//// Clean up the container after the test is complete
	t.Cleanup(func() {
		if terr := dbContainer.Terminate(ctx); terr != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// get a free random host port for rabbitmq `Tcp Port`
	hostPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port("5672/tcp"),
	)
	if err != nil {
		return nil, err
	}

	// https://github.com/michaelklishin/rabbit-hole/issues/74
	// get a free random host port for rabbitmq UI `Http Port`
	uiHttpPort, err := dbContainer.MappedPort(
		ctx,
		nat.Port("15672/tcp"),
	)
	if err != nil {
		return nil, err
	}

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	option := &config.RabbitmqHostOptions{
		UserName:    "guest",
		Password:    "guest",
		HostName:    host,
		VirtualHost: "/",
		Port:        hostPort.Int(),
		HttpPort:    uiHttpPort.Int(),
	}

	return option, nil
}
