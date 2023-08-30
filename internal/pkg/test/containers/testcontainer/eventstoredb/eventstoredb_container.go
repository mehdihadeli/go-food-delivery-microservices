package eventstoredb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"
)

type eventstoredbTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.EventstoreDBContainerOptions
}

func NewEventstoreDBTestContainers() contracts.EventstoreDBContainer {
	return &eventstoredbTestContainers{
		defaultOptions: &contracts.EventstoreDBContainerOptions{
			Ports:   []string{"2113/tcp", "1113/tcp"},
			Host:    "localhost",
			TcpPort: 1113,
			// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
			HttpPort:  2113,
			Tag:       "latest",
			ImageName: "eventstore/eventstore",
			Name:      "eventstoredb-testcontainers",
		},
	}
}

func (g *eventstoredbTestContainers) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.EventstoreDBContainerOptions,
) (*config.EventStoreDbOptions, error) {
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

	// get a free random host port for http and grpc port for eventstoredb
	httpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[0]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HttpPort = httpPort.Int()
	t.Logf("eventstoredb http and grpc port is: %d", httpPort.Int())

	// get a free random host port for tcp port eventstoredb
	tcpPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Ports[1]))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.TcpPort = tcpPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	option := &config.EventStoreDbOptions{
		Host:     host,
		TcpPort:  g.defaultOptions.TcpPort,
		HttpPort: g.defaultOptions.HttpPort,
	}

	return option, nil
}

func (g *eventstoredbTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.EventstoreDBContainerOptions,
) (*esdb.Client, error) {
	eventstoredbOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}
	return eventstroredb.NewEventStoreDB(eventstoredbOptions)
}

func (g *eventstoredbTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
}

func (g *eventstoredbTestContainers) getRunOptions(
	opts ...*contracts.EventstoreDBContainerOptions,
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
		if option.Tag != "" {
			g.defaultOptions.Tag = option.Tag
		}
	}

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: g.defaultOptions.Ports,
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Ports[0])).WithPollInterval(2 * time.Second),
		Hostname:     g.defaultOptions.Host,
		SkipReaper:   true,
		// we use `EVENTSTORE_IN_MEM` for use eventstoredb in-memory mode in tests
		Env: map[string]string{
			"EVENTSTORE_START_STANDARD_PROJECTIONS": "false",
			"EVENTSTORE_INSECURE":                   "true",
			"EVENTSTORE_ENABLE_EXTERNAL_TCP":        "true",
			"EVENTSTORE_ENABLE_ATOM_PUB_OVER_HTTP":  "true",
			"EVENTSTORE_MEM_DB":                     "true",
		},
	}

	return containerReq
}
