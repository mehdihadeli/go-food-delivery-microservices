package orders

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

type ordersServiceConfigurations struct {
	infra            contracts2.InfrastructureConfigurations
	ordersMetrics    *contracts2.OrdersMetrics
	ordersEchoServer customEcho.EchoHttpServer
	ordersGrpcServer grpcServer.GrpcServer
	ordersBus        bus.Bus
	esdbWorker       eventstroredb.EsdbSubscriptionAllWorker
}
