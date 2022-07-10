package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/kafka_messages"
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToGrpcMessage(product *models.Product) *kafka_messages.Product {
	return &kafka_messages.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ProductFromGrpcMessage(product *kafka_messages.Product) (*models.Product, error) {

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

//func ProductToProductDto(product *models.Product) *dto.ProductDto {
//	return &dto.ProductDto{
//		ProductID:   product.ProductID,
//		Name:        product.Name,
//		Description: product.Description,
//		Price:       product.Price,
//		CreatedAt:   product.CreatedAt,
//		UpdatedAt:   product.UpdatedAt,
//	}
//}
//
//func ProductsToProductsDto(products []*models.Product) []*dto.ProductDto {
//	productDtos := make([]*dto.ProductDto, 0, len(products))
//	for _, product := range products {
//		productDtos = append(productDtos, ProductToProductDto(product))
//	}
//
//	return productDtos
//}
