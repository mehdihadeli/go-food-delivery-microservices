package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	return nil
}

