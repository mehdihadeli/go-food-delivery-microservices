package contracts

import (
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
)

type CatalogsServiceConfigurations struct {
	Infra              *InfrastructureConfigurations
	CatalogsMetrics    *CatalogsMetrics
	CatalogsEchoServer customEcho.EchoHttpServer
	CatalogsGrpcServer grpcServer.GrpcServer
	CatalogsBus        bus.Bus
}
