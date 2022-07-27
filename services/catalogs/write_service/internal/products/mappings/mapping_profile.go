package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/kafka_messages"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(product *models.Product) *kafka_messages.Product {
		return &kafka_messages.Product{
			ProductID:   product.ProductID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(product *kafka_messages.Product) *models.Product {
		proUUID, err := uuid.FromString(product.GetProductID())
		if err != nil {
			return nil
		}

		return &models.Product{
			ProductID:   proUUID,
			Name:        product.GetName(),
			Description: product.GetDescription(),
			Price:       product.GetPrice(),
			CreatedAt:   product.GetCreatedAt().AsTime(),
			UpdatedAt:   product.GetUpdatedAt().AsTime(),
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(product *models.Product) *productsService.Product {
		return &productsService.Product{
			ProductID:   product.ProductID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		}
	})

	return nil
}
