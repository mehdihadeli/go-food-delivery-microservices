package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 60 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type mongoTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.MongoContainerOptions
	logger         logger.Logger
}

func NewMongoTestContainers(l logger.Logger) contracts.MongoContainer {
	return &mongoTestContainers{
		defaultOptions: &contracts.MongoContainerOptions{
			Database:  "test_db",
			Port:      "27017/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "mongo",
			Name:      "mongo-testcontainer",
		},
		logger: l,
	}
}

func (g *mongoTestContainers) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongodb.MongoDbOptions, error) {
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
		return g.CreatingContainerOptions(context.Background(), t, options...)
	}

	g.container = dbContainer

	option := &mongodb.MongoDbOptions{
		User:     g.defaultOptions.UserName,
		Password: g.defaultOptions.Password,
		UseAuth:  false,
		Host:     host,
		Port:     g.defaultOptions.HostPort,
		Database: g.defaultOptions.Database,
	}

	return option, nil
}

func (g *mongoTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongo.Client, error) {
	mongoOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := mongodb.NewMongoDB(mongoOptions)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (g *mongoTestContainers) Cleanup(ctx context.Context) error {
	if err := g.container.Terminate(ctx); err != nil {
		return errors.WrapIf(err, "failed to terminate container: %s")
	}

	return nil
}

func (g *mongoTestContainers) getRunOptions(
	opts ...*contracts.MongoContainerOptions,
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

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Port)).WithPollInterval(2 * time.Second),
		Hostname:     g.defaultOptions.Host,
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": g.defaultOptions.UserName,
			"MONGO_INITDB_ROOT_PASSWORD": g.defaultOptions.Password,
		},
	}

	return containerReq
}

func isConnectable(ctx context.Context, logger logger.Logger, mongoOptions *contracts.MongoContainerOptions) bool {
	uriAddress := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		mongoOptions.UserName,
		mongoOptions.Password,
		mongoOptions.Host,
		mongoOptions.HostPort,
	)
	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)
	opt = opt.SetAuth(options.Credential{Username: mongoOptions.UserName, Password: mongoOptions.Password})

	mongoClient, err := mongo.Connect(ctx, opt)

	defer mongoClient.Disconnect(ctx)

	if err != nil {
		logError(logger, mongoOptions.Host, mongoOptions.HostPort)

		return false
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		logError(logger, mongoOptions.Host, mongoOptions.HostPort)

		return false
	}
	logger.Infof(
		"Opened mongodb connection on host: %s:%d", mongoOptions.Host, mongoOptions.HostPort)

	return true
}

func logError(logger logger.Logger, host string, hostPort int) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf(
		"Error in creating mongodb connection with %s:%d", host, hostPort,
	)
}
