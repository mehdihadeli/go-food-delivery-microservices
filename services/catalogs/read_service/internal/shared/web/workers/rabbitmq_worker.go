package workers

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func NewRabbitMQWorker(ctx context.Context, infra contracts.InfrastructureConfiguration) web.Worker {
	return web.NewBackgroundWorker(func(ctx context.Context) error {
		err := infra.GetRabbitMQBus().Start(ctx)
		if err != nil {
			infra.GetLog().Errorf("[RabbitMQWorkerWorker.Start] error in the starting rabbitmq worker: {%v}", err)
			return err
		}
		return nil
	}, func(ctx context.Context) error {
		return infra.GetRabbitMQBus().Stop(ctx)
	})
}
