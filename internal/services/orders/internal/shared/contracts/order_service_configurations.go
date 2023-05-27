package contracts

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
)

type OrderServiceConfigurations struct {
	InfrastructureConfigurations *InfrastructureConfigurations
	OrdersMetrics                *OrdersMetrics
	OrdersEchoServer             customEcho.EchoHttpServer
	OrdersGrpcServer             grpcServer.GrpcServer
	OrdersBus                    bus.Bus
	OrdersSubscriptionAllWorker  eventstroredb.EsdbSubscriptionAllWorker
}
