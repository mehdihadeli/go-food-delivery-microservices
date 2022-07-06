package mappers

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/grpc/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
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

func ProductToProductDto(product *models.Product) *dto.ProductDto {
	return &dto.ProductDto{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ProductsToProductsDto(products []*models.Product) []*dto.ProductDto {
	productDtos := make([]*dto.ProductDto, 0, len(products))
	for _, product := range products {
		productDtos = append(productDtos, ProductToProductDto(product))
	}

	return productDtos
}

func ListResultToListResultDto[TModel any, TDto any](listResult *utils.ListResult[TModel], m func([]*TModel) []*TDto) *utils.ListResult[TDto] {
	return &utils.ListResult[TDto]{
		Items:      m(listResult.Items),
		Size:       listResult.Size,
		Page:       listResult.Page,
		TotalItems: listResult.TotalItems,
		TotalPage:  listResult.TotalPage,
	}
}
