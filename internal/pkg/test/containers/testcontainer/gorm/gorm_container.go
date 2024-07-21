package gorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/contracts"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// https://github.com/testcontainers/testcontainers-go/issues/1359
// https://github.com/testcontainers/testcontainers-go/issues/1249
type gormTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.PostgresContainerOptions
	logger         logger.Logger
}

func NewGormTestContainers(l logger.Logger) contracts.GormContainer {
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
		logger: l,
	}
}

func (g *gormTestContainers) PopulateContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gormPostgres.GormOptions, error) {
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
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

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := dbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		time.Sleep(time.Second * 1)
	})

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

	isConnectable := isConnectable(ctx, g.logger, g.defaultOptions)
	if !isConnectable {
		return g.PopulateContainerOptions(context.Background(), t, options...)
	}

	g.container = dbContainer

	gormOptions := &gormPostgres.GormOptions{
		Port:     g.defaultOptions.HostPort,
		Host:     host,
		Password: g.defaultOptions.Password,
		DBName:   g.defaultOptions.Database,
		SSLMode:  false,
		User:     g.defaultOptions.UserName,
	}
	return gormOptions, nil
}

func (g *gormTestContainers) Start(
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

func (g *gormTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}
	return nil
}

func (g *gormTestContainers) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) testcontainers.ContainerRequest {
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

	strategies := []wait.Strategy{wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)}
	deadline := 120 * time.Second

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor:   wait.ForAll(strategies...).WithDeadline(deadline),
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_DB":       g.defaultOptions.Database,
			"POSTGRES_PASSWORD": g.defaultOptions.Password,
			"POSTGRES_USER":     g.defaultOptions.UserName,
		},
	}

	return containerReq
}

func isConnectable(
	ctx context.Context,
	logger logger.Logger,
	postgresOptions *contracts.PostgresContainerOptions,
) bool {
	orm, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%d/postgres?sslmode=disable",
				postgresOptions.UserName,
				postgresOptions.Password,
				postgresOptions.Host,
				postgresOptions.HostPort,
			),
		),
		&gorm.Config{
			PrepareStmt:              true,
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		},
	)
	if err != nil {
		logError(logger, postgresOptions.Host, postgresOptions.HostPort)

		return false
	}

	db, err := orm.DB()
	if err != nil {
		logError(logger, postgresOptions.Host, postgresOptions.HostPort)

		return false
	}

	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		logError(logger, postgresOptions.Host, postgresOptions.HostPort)

		return false
	}

	logger.Infof(
		"Opened postgres connection on host: %s:%d", postgresOptions.Host, postgresOptions.HostPort)

	return true
}

func logError(logger logger.Logger, host string, hostPort int) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf("Error in creating postgres connection with %s:%d", host, hostPort)
}
