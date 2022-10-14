package contracts

import (
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
)

type CatalogsServiceConfigurations struct {
	Infra              *InfrastructureConfigurations
	CatalogsMetrics    *CatalogsMetrics
	CatalogsEchoServer customEcho.EchoHttpServer
	CatalogsGrpcServer grpcServer.GrpcServer
	CatalogsBus        bus.Bus
}
