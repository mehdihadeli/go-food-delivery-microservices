package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
