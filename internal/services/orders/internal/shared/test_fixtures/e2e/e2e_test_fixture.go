package e2e

import (
	"context"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	webWoker "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

type E2ETestFixture struct {
	Echo          *echo.Echo
	V1            *V1Groups
	GrpcServer    grpcServer.GrpcServer
	HttpServer    *httptest.Server
	workersRunner *webWoker.WorkersRunner
	Bus           bus.Bus
	OrdersMetrics *contracts.OrdersMetrics
	Ctx           context.Context
	cancel        context.CancelFunc
	Cleanup       func()
}

type V1Groups struct {
	OrdersGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	//cfg, _ := config.NewConfig(constants.Test)
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//
	//echo := echo.New()
	//
	//v1Group := echo.Group("/api/v1")
	//ordersV1 := v1Group.Group("/orders")
	//
	//v1Groups := &V1Groups{OrdersGroup: ordersV1}
	//
	//// this should not be in integration_events test because of cyclic dependencies
	//err = mediatr.ConfigOrdersMediator(infrastructures)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//err = mappings.ConfigureOrdersMappings()
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//mq, err := rabbitmq.ConfigOrdersRabbitMQ(ctx, cfg.RabbitMQ, infrastructures)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//subscriptionAllWorker, err := subscriptionAll.ConfigOrdersSubscriptionAllWorker(
	//	infrastructures,
	//	mq,
	//)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//ordersMetrics, err := metrics.ConfigOrdersMetrics(cfg, infrastructures.Metrics)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//grpcServer := grpcServer.NewGrpcServer(
	//	cfg.GRPC,
	//	defaultLogger2.Logger,
	//	cfg.ServiceName,
	//	infrastructures.Metrics,
	//)
	//httpServer := httptest.NewServer(echo)
	//
	//workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
	//	workers.NewRabbitMQWorker(
	//		infrastructures.Log,
	//		mq,
	//	), workers.NewEventStoreDBWorker(infrastructures.Log, infrastructures.Cfg, subscriptionAllWorker),
	//})
	//
	//return &E2ETestFixture{
	//	Cleanup: func() {
	//		workersRunner.Stop(ctx)
	//		cancel()
	//		cleanup()
	//		grpcServer.GracefulShutdown()
	//		echo.Shutdown(ctx)
	//		httpServer.Close()
	//	},
	//	InfrastructureConfigurations: infrastructures,
	//	Echo:                         echo,
	//	V1:                           v1Groups,
	//	Bus:                          mq,
	//	OrdersMetrics:                ordersMetrics,
	//	GrpcServer:                   grpcServer,
	//	HttpServer:                   httpServer,
	//	workersRunner:                workersRunner,
	//	Ctx:                          ctx,
	//	cancel:                       cancel,
	//}

	return &E2ETestFixture{}
}

func (e *E2ETestFixture) Run() {
	workersErr := e.workersRunner.Start(e.Ctx)
	go func() {
		for {
			select {
			case _ = <-workersErr:
				e.cancel()
				return
			}
		}
	}()
}
