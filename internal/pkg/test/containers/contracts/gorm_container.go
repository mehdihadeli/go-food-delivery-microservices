package contracts

import (
	"context"
	"testing"

	"gorm.io/gorm"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
)

type PostgresContainerOptions struct {
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

type GormContainer interface {
	Start(ctx context.Context, t *testing.T, options ...*PostgresContainerOptions) (*gorm.DB, error)
	CreatingContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*gormPostgres.GormOptions, error)
	Cleanup(ctx context.Context) error
}
