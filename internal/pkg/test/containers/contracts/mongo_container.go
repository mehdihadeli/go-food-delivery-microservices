package contracts

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mongodb"
)

type MongoContainerOptions struct {
	Database  string
	Host      string
	Port      string
	HostPort  int
	UserName  string
	Password  string
	ImageName string
	Name      string
	Tag       string
}

type MongoContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*MongoContainerOptions,
	) (*mongodb.MongoDbOptions, error)
	Cleanup(ctx context.Context) error
}
