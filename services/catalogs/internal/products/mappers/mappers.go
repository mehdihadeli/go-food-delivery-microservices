package mappers

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/proto/product_kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/proto/product_service"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToGrpcMessage(product *models.Product) *product_kafka_messages.Product {
	return &product_kafka_messages.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ProductFromGrpcMessage(product *product_kafka_messages.Product) (*models.Product, error) {

	proUUID, err := uuid.FromString(product.GetProductID())
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ProductID:   proUUID,
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		CreatedAt:   product.GetCreatedAt().AsTime(),
		UpdatedAt:   product.GetUpdatedAt().AsTime(),
	}, nil
}

func WriterProductToGrpc(product *models.Product) *product_service.Product {
	return &product_service.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}
