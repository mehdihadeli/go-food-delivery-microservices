package orders

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

type ordersServiceConfigurations struct {
	infra            contracts2.InfrastructureConfigurations
	ordersMetrics    contracts2.OrdersMetrics
	ordersEchoServer customEcho.EchoHttpServer
	ordersGrpcServer grpcServer.GrpcServer
	ordersBus        bus.Bus
	esdbWorker       eventstroredb.EsdbSubscriptionAllWorker
}

func (c *ordersServiceConfigurations) OrdersMetrics() contracts2.OrdersMetrics {
	return c.ordersMetrics
}

func (c *ordersServiceConfigurations) OrdersGrpcServer() grpcServer.GrpcServer {
	return c.ordersGrpcServer
}

func (c *ordersServiceConfigurations) OrdersBus() bus.Bus {
	return c.ordersBus
}

func (c *ordersServiceConfigurations) InfrastructureConfigurations() contracts2.InfrastructureConfigurations {
	return c.infra
}

func (c *ordersServiceConfigurations) OrdersEchoServer() customEcho.EchoHttpServer {
	return c.ordersEchoServer
}

func (c *ordersServiceConfigurations) OrdersSubscriptionAllWorker() eventstroredb.EsdbSubscriptionAllWorker {
	return c.esdbWorker
}
