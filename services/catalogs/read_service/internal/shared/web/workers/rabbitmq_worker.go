package workers

import (
	"context"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

func NewRabbitMQWorkerWorker(infra *infrastructure.InfrastructureConfigurations) web.Worker {
	rabbitMQBus := rabbitmqBus.NewRabbitMQBus(infra.Log, infra.Consumers...)

	return web.NewBackgroundWorker(func(ctx context.Context) error {
		err := rabbitMQBus.Start(ctx)
		if err != nil {
			infra.Log.Errorf("[RabbitMQWorkerWorker.Start] error in the starting rabbitmq worker: {%v}", err)
			return err
		}
		return nil
	}, func(ctx context.Context) error {
		return rabbitMQBus.Stop(ctx)
	})
}
