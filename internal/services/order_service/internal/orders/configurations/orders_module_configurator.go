package configurations

import (
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/store"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	googleGrpc "google.golang.org/grpc"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/grpc"
	ordersservice "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/grpc/genproto"
)

type OrdersModuleConfigurator struct {
	contracts2.Application
}

func NewOrdersModuleConfigurator(
	app contracts2.Application,
) *OrdersModuleConfigurator {
	return &OrdersModuleConfigurator{
		Application: app,
	}
}

func (c *OrdersModuleConfigurator) ConfigureOrdersModule() {
	c.ResolveFunc(
		func(logger logger.Logger,
			server customEcho.EchoHttpServer,
			orderRepository repositories.OrderMongoRepository,
			orderAggregateStore store.AggregateStore[*aggregate.Order],
			tracer tracing.AppTracer,
		) error {
			// Config Orders Mappings
			err := mappings.ConfigureOrdersMappings()
			if err != nil {
				return err
			}

			// Config Orders Mediators
			err = mediatr.ConfigOrdersMediator(logger, orderRepository, orderAggregateStore, tracer)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (c *OrdersModuleConfigurator) MapOrdersEndpoints() {
	// Config Orders Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"order-routes"`,
	)

	// Config Orders Grpc Endpoints
	c.ResolveFunc(
		func(ordersGrpcServer grpcServer.GrpcServer, ordersMetrics *contracts.OrdersMetrics, logger logger.Logger, validator *validator.Validate) error {
			orderGrpcService := grpc.NewOrderGrpcService(logger, validator, ordersMetrics)
			ordersGrpcServer.GrpcServiceBuilder().RegisterRoutes(func(server *googleGrpc.Server) {
				ordersservice.RegisterOrdersServiceServer(server, orderGrpcService)
			})
			return nil
		},
	)
}
