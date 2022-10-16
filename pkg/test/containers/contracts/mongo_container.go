package contracts

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
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
	Start(ctx context.Context, t *testing.T, options ...*MongoContainerOptions) (*mongo.Client, error)
	Cleanup(ctx context.Context) error
}
