package mappings

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtoV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"
	productsService "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureProductsMappings() error {
	err := mapper.CreateMap[*models.Product, *dtoV1.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dtoV1.ProductDto, *models.Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*datamodel.ProductDataModel, *models.Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *datamodel.ProductDataModel]()
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*dtoV1.ProductDto, *productsService.Product](
		func(product *dtoV1.ProductDto) *productsService.Product {
			if product == nil {
				return nil
			}
			return &productsService.Product{
				ProductId:   product.Id.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}
		},
	)
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap(
		func(product *models.Product) *productsService.Product {
			return &productsService.Product{
				ProductId:   product.Id.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}
		},
	)

	return nil
}
