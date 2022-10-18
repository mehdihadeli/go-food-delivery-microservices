package uow

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres/uow"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"gorm.io/gorm"
)

type catalogsUnitOfWork struct {
	productRepository contracts.ProductRepository
	data.UnitOfWork
	db *gorm.DB
}

func NewCatalogsUnitOfWork(logger logger.Logger, db *gorm.DB, productRepository contracts.ProductRepository) contracts.CatalogsUnitOfWorks {
	return &catalogsUnitOfWork{UnitOfWork: uow.NewGormUnitOfWork(db, logger), productRepository: productRepository}
}

func (c *catalogsUnitOfWork) Products() contracts.ProductRepository {
	return c.productRepository
}
