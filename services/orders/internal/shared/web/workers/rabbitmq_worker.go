package workers

import (
	"context"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func NewRabbitMQWorker(ctx context.Context, infra *infrastructure.InfrastructureConfiguration) web.Worker {
	bus, _ := rabbitmqBus.NewRabbitMQBus(
		ctx,
		infra.Cfg.RabbitMQ,
		func(builder configurations.RabbitMQConfigurationBuilder) {
			rabbitmq.ConfigRabbitMQ(builder, infra)
		},
		infra.EventSerializer,
		infra.Log)

	infra.RabbitMQBus = bus
	infra.Producer = bus

	return web.NewBackgroundWorker(func(ctx context.Context) error {
		err := bus.Start(ctx)
		if err != nil {
			infra.Log.Errorf("[RabbitMQWorkerWorker.Start] error in the starting rabbitmq worker: {%v}", err)
			return err
		}
		return nil
	}, func(ctx context.Context) error {
		return bus.Stop(ctx)
	})
}
