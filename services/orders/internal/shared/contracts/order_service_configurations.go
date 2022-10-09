package contracts

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
)

type OrderServiceConfigurations interface {
	InfrastructureConfigurations() InfrastructureConfigurations
	OrdersMetrics() OrdersMetrics
	OrdersEchoServer() customEcho.EchoHttpServer
	OrdersGrpcServer() grpcServer.GrpcServer
	OrdersBus() bus.Bus
	OrdersSubscriptionAllWorker() eventstroredb.EsdbSubscriptionAllWorker
}
