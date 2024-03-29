package contracts

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
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
	Start(
		ctx context.Context,
		t *testing.T,
		options ...*MongoContainerOptions,
	) (*mongo.Client, error)
	CreatingContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*MongoContainerOptions,
	) (*mongodb.MongoDbOptions, error)
	Cleanup(ctx context.Context) error
}
