package gorm

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/contracts"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
	"testing"
)

type gormTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.PostgresContainerOptions
}

func NewGormTestContainers() *gormTestContainers {
	return &gormTestContainers{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-testcontainer",
		},
	}
}

func (g *gormTestContainers) Start(ctx context.Context, t *testing.T, options ...*contracts.PostgresContainerOptions) (*gorm.DB, error) {
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
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Port))
	if err != nil {
		return nil, err
	}
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

	db, err := gormPostgres.NewGorm(&gormPostgres.GormConfig{
		Port:     g.defaultOptions.HostPort,
		Host:     host,
		Password: g.defaultOptions.Password,
		DBName:   g.defaultOptions.Database,
		SSLMode:  false,
		User:     g.defaultOptions.UserName,
	})

	return db, nil
}

func (g *gormTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
}

func (g *gormTestContainers) getRunOptions(opts ...*contracts.PostgresContainerOptions) testcontainers.ContainerRequest {
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

	//hostFreePort, err := freeport.GetFreePort()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//g.defaultOptions.HostPort = hostFreePort

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Port)),
		Hostname:     g.defaultOptions.Host,
		Env: map[string]string{
			"POSTGRES_DB":       g.defaultOptions.Database,
			"POSTGRES_PASSWORD": g.defaultOptions.Password,
			"POSTGRES_USER":     g.defaultOptions.UserName,
		},
	}

	return containerReq
}
