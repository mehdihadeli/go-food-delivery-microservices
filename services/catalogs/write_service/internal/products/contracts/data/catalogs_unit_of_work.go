package data

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"

type CatalogsUnitOfWorks interface {
	data.UnitOfWork
	Products() ProductRepository
}
