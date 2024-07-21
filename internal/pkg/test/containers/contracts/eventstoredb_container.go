package contracts

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/eventstroredb/config"
)

type EventstoreDBContainerOptions struct {
	Host    string
	Ports   []string
	TcpPort int
	// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
	HttpPort  int
	ImageName string
	Name      string
	Tag       string
}

type EventstoreDBContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*EventstoreDBContainerOptions,
	) (*config.EventStoreDbOptions, error)

	Cleanup(ctx context.Context) error
}
