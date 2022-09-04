package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
)

//https://echo.labstack.com/guide/response/

type GetProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
