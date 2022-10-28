package uow

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres/uow"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	data2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
	"gorm.io/gorm"
)

type catalogsUnitOfWork struct {
	productRepository data2.ProductRepository
	data.UnitOfWork
	db *gorm.DB
}

func NewCatalogsUnitOfWork(logger logger.Logger, db *gorm.DB, productRepository data2.ProductRepository) data2.CatalogsUnitOfWorks {
	return &catalogsUnitOfWork{UnitOfWork: uow.NewGormUnitOfWork(db, logger), productRepository: productRepository}
}

func (c *catalogsUnitOfWork) Products() data2.ProductRepository {
	return c.productRepository
}
