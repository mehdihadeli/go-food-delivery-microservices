package uow

import (
	data2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
)

type catalogContext struct {
	productRepository data2.ProductRepository
}

func (c *catalogContext) Products() data2.ProductRepository {
	return c.productRepository
}
