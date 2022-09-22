package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
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
