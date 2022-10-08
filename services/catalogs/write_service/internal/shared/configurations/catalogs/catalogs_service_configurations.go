package catalogs

import (
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
)

type catalogsServiceConfigurations struct {
	infra              contracts.InfrastructureConfigurations
	catalogsMetrics    contracts.CatalogsMetrics
	catalogsEchoServer customEcho.EchoHttpServer
	catalogsGrpcServer grpcServer.GrpcServer
	catalogsBus        bus.Bus
}

func (c *catalogsServiceConfigurations) InfrastructureConfigurations() contracts.InfrastructureConfigurations {
	return c.infra
}

func (c *catalogsServiceConfigurations) CatalogsMetrics() contracts.CatalogsMetrics {
	return c.catalogsMetrics
}

func (c *catalogsServiceConfigurations) CatalogsGrpcServer() grpcServer.GrpcServer {
	return c.catalogsGrpcServer
}

func (c *catalogsServiceConfigurations) CatalogsEchoServer() customEcho.EchoHttpServer {
	return c.catalogsEchoServer
}
func (c *catalogsServiceConfigurations) CatalogsBus() bus.Bus {
	return c.catalogsBus
}
