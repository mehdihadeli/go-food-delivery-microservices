package gorm

import (
	"context"
	"log"
	"testing"

	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/contracts"

	"emperror.dev/errors"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"gorm.io/gorm"
)

type gnoMockGormContainer struct {
	container      *gnomock.Container
	defaultOptions *contracts.PostgresContainerOptions
}

func NewGnoMockGormContainer() contracts.GormContainer {
	return &gnoMockGormContainer{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432",
			Host:      "localhost",
			UserName:  "genomocktest",
			Password:  "genomocktest",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-genomock-container",
		},
	}
}

func (g *gnoMockGormContainer) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gormPostgres.GormOptions, error) {
	// https://github.com/orlangure/gnomock
	gnomock.WithContext(ctx)
	runOption := g.getRunOptions(options...)
	container, err := gnomock.Start(runOption)
	if container == nil || err != nil {
		log.Fatal(errors.New("error creating postgres container"))
		return nil, err
	}

	g.container = container
	g.defaultOptions.HostPort = container.DefaultPort()

	t.Cleanup(func() { _ = gnomock.Stop(container) })

	gormContainerOptions := &gormPostgres.GormOptions{
		Port:     g.defaultOptions.HostPort,
		Host:     container.Host,
		Password: g.defaultOptions.Password,
		DBName:   g.defaultOptions.Database,
		SSLMode:  false,
		User:     g.defaultOptions.UserName,
	}

	return gormContainerOptions, nil
}

func (g *gnoMockGormContainer) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gorm.DB, error) {
	gormOptions, err := g.PopulateContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := gormPostgres.NewGorm(gormOptions)

	return db, nil
}

func (g *gnoMockGormContainer) Cleanup(ctx context.Context) error {
	return gnomock.Stop(g.container)
}

func (g *gnoMockGormContainer) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) gnomock.Preset {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if option.Port != "" {
			g.defaultOptions.Port = option.Port
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

	p := postgres.Preset(
		postgres.WithDatabase(g.defaultOptions.Database),
		postgres.WithUser(g.defaultOptions.UserName, g.defaultOptions.Password),
		postgres.WithVersion(g.defaultOptions.Tag),
	)

	return p
}
