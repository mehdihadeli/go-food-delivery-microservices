package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/dto"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
