package grpc

import (
	"context"

	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsGrpc(ctx context.Context, builder *grpcServer.GrpcServiceBuilder, infra *contracts.InfrastructureConfigurations) {
}
