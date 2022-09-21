package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/events/integration"
	UpdatingProductIntegration "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/events/integration"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dto.ProductDto, *models.Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*dto.ProductDto, *productsService.Product](func(product *dto.ProductDto) *productsService.Product {
		if product == nil {
			return nil
		}
		return &productsService.Product{
			ProductId:   product.ProductId.String(),
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

	err = mapper.CreateCustomMap(func(product *models.Product) *integration.ProductCreated {
		return &integration.ProductCreated{
			ProductId:   product.ProductId.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			Message:     types.NewMessage(uuid.NewV4().String()),
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(product *models.Product) *UpdatingProductIntegration.ProductUpdated {
		return &UpdatingProductIntegration.ProductUpdated{
			ProductId:   product.ProductId.String(),
			UpdatedAt:   product.UpdatedAt,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Message:     types.NewMessage(uuid.NewV4().String()),
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(productUpdated *UpdatingProductIntegration.ProductUpdated) *models.Product {
		proUUID, err := uuid.FromString(productUpdated.ProductId)
		if err != nil {
			return nil
		}

		return &models.Product{
			ProductId:   proUUID,
			Name:        productUpdated.Name,
			Description: productUpdated.Description,
			Price:       productUpdated.Price,
			UpdatedAt:   productUpdated.UpdatedAt,
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(func(product *models.Product) *productsService.Product {
		return &productsService.Product{
			ProductId:   product.ProductId.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		}
	})

	return nil
}
