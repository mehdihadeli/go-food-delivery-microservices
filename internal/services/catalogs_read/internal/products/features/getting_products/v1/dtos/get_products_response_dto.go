package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/dto"
)

type GetProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
