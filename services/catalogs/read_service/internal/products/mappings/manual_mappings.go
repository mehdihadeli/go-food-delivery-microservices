package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/grpc/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToGrpcMessage(product *models.Product) *kafka_messages.Product {
	return &kafka_messages.Product{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}