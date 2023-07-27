package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/dto"
)

type GetProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
