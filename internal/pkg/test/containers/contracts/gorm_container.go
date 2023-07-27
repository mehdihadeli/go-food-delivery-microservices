package contracts

import (
	"context"
	"gorm.io/gorm"
	"testing"
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
	Cleanup(ctx context.Context) error
}
