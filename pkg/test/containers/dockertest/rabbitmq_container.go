package dockertest

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/contracts"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"strconv"
	"testing"
)

type rabbitmqDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.RabbitMQContainerOptions
}

func NewRabbitMQDockerTest() contracts.RabbitMQContainer {
	return &rabbitmqDockerTest{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672", "15672"},
			Host:        "localhost",
			VirtualHost: "",
			UserName:    "dockertest",
			Password:    "dockertest",
			Tag:         "3-management",
			ImageName:   "rabbitmq",
			Name:        "rabbitmq-dockertest",
		},
	}
}

func (g *rabbitmqDockerTest) Start(ctx context.Context, t *testing.T, rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc, options ...*contracts.RabbitMQContainerOptions) (bus.Bus, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOptions := g.getRunOptions(options...)

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(runOptions, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource (RabbitMQ Container): %s", err)
	}

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application in the container might not be ready to accept connections yet

	g.resource = resource
	i, err := strconv.Atoi(resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Ports[0]))) //5672
	g.defaultOptions.HostPort = i

	t.Cleanup(func() { _ = resource.Close() })

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = resource.Close()
				return
			}
		}
	}()

	var mqBus bus.Bus
	if err = pool.Retry(func() error {
		mqBus, err = bus2.NewRabbitMQBus(
			ctx,
			&config.RabbitMQConfig{
				RabbitMqHostOptions: &config.RabbitMqHostOptions{
					UserName:    g.defaultOptions.UserName,
					Password:    g.defaultOptions.Password,
					HostName:    g.defaultOptions.Host,
					VirtualHost: g.defaultOptions.VirtualHost,
					Port:        g.defaultOptions.HostPort,
				},
			},
			rabbitmqBuilderFunc,
			json.NewJsonEventSerializer(),
			defaultLogger.Logger)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
		return nil, err
	}

	return mqBus, nil
}

func (g *rabbitmqDockerTest) Cleanup(ctx context.Context) error {
	return g.resource.Close()
}

func (g *rabbitmqDockerTest) getRunOptions(opts ...*contracts.RabbitMQContainerOptions) *dockertest.RunOptions {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if len(option.Ports) > 0 {
			g.defaultOptions.Ports = option.Ports
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

	runOptions := &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"RABBITMQ_DEFAULT_USER=" + g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS=" + g.defaultOptions.Password,
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: g.defaultOptions.Ports,
	}

	return runOptions
}
