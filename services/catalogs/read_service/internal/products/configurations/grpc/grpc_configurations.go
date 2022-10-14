package grpc

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra *contracts.InfrastructureConfigurations) {
}
