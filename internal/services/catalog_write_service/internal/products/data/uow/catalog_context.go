package uow

import (
	data2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts"
)

type catalogContext struct {
	productRepository data2.ProductRepository
}

func (c *catalogContext) Products() data2.ProductRepository {
	return c.productRepository
}
