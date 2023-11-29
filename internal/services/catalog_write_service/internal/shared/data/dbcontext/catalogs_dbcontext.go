package dbcontext

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/helpers"
	datamodel "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CatalogContextActionFunc func(ctx context.Context, catalogContext *CatalogsGormDBContext) error

type CatalogsGormDBContext struct {
	*gorm.DB
	logger logger.Logger
}

func NewCatalogsDBContext(
	db *gorm.DB,
	log logger.Logger,
) *CatalogsGormDBContext {
	c := &CatalogsGormDBContext{DB: db, logger: log}

	return c
}

// WithTx creates a transactional DBContext with getting tx-gorm from the ctx. This will throw an error if the transaction does not exist.
func (c *CatalogsGormDBContext) WithTx(
	ctx context.Context,
) (*CatalogsGormDBContext, error) {
	tx, err := helpers.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return NewCatalogsDBContext(tx, c.logger), nil
}

// WithTxIfExists creates a transactional DBContext with getting tx-gorm from the ctx. not throw an error if the transaction is not existing and returns an existing database.
func (c *CatalogsGormDBContext) WithTxIfExists(
	ctx context.Context,
) *CatalogsGormDBContext {
	tx := helpers.GetTxFromContextIfExists(ctx)
	if tx == nil {
		return c
	}

	return NewCatalogsDBContext(tx, c.logger)
}

func (c *CatalogsGormDBContext) RunInTx(
	ctx context.Context,
	action CatalogContextActionFunc,
) error {
	// https://gorm.io/docs/transactions.html#Transaction
	tx := c.WithContext(ctx).Begin()

	c.logger.Info("beginning database transaction")

	gormContext := helpers.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, _ := r.(error); err != nil {
				c.logger.Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				c.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	err := action(ctx, c)
	if err != nil {
		c.logger.Error("rolling back transaction")
		tx.WithContext(ctx).Rollback()

		return err
	}

	c.logger.Info("committing transaction")

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		c.logger.Errorf("transaction commit error: ", err)
	}

	return err
}

// Extensions for shared and reusable methods. for complex and none-reusable queries we use DBContext directly

func (c *CatalogsGormDBContext) FindProductByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Product, error) {
	var productDatas []*datamodel.ProductDataModel
	var productData *datamodel.ProductDataModel

	s := c.DB.WithContext(ctx).Find(&productDatas).Error
	fmt.Println(s)

	// https://gorm.io/docs/query.html#Retrieving-objects-with-primary-key
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	// https://gorm.io/docs/query.html#Inline-Condition
	// https://gorm.io/docs/advanced_query.html
	result := c.WithContext(ctx).First(&productData, id)
	if result.Error != nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			result.Error,
			fmt.Sprintf(
				"product with id `%s` not found in the database",
				id.String(),
			),
		)
	}

	c.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	product, err := mapper.Map[*models.Product](productData)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping Product",
		)
	}

	return product, nil
}

// DeleteProductByID delete the product inner a tx if exists
func (c *CatalogsGormDBContext) DeleteProductByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	dbContext := c.WithTxIfExists(ctx)

	product, err := dbContext.FindProductByID(ctx, id)
	if err != nil {
		return customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"product with id `%s` not found in the database",
				id.String(),
			),
		)
	}

	productDataModel, err := mapper.Map[*datamodel.ProductDataModel](product)
	if err != nil {
		return customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping ProductDataModel",
		)
	}

	result := dbContext.WithContext(ctx).Delete(productDataModel, id)
	if result.Error != nil {
		return customErrors.NewInternalServerErrorWrap(
			result.Error,
			fmt.Sprintf(
				"error in deleting product with id `%s` in the database",
				id.String(),
			),
		)
	}

	c.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return nil
}

// AddProduct add the product inner a tx if exists
func (c *CatalogsGormDBContext) AddProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	dbContext := c.WithTxIfExists(ctx)

	productDataModel, err := mapper.Map[*datamodel.ProductDataModel](product)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping ProductDataModel",
		)
	}

	// https://gorm.io/docs/create.html
	result := dbContext.WithContext(ctx).Create(productDataModel)
	if result.Error != nil {
		return nil, customErrors.NewConflictErrorWrap(
			result.Error,
			"product already exists",
		)
	}

	c.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	product, err = mapper.Map[*models.Product](productDataModel)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping Product",
		)
	}

	return product, err
}

// UpdateProduct update the product inner a tx if exists
func (c *CatalogsGormDBContext) UpdateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	dbContext := c.WithTxIfExists(ctx)

	productDataModel, err := mapper.Map[*datamodel.ProductDataModel](product)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping ProductDataModel",
		)
	}

	// https://gorm.io/docs/update.html
	result := dbContext.WithContext(ctx).Updates(productDataModel)
	if result.Error != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			result.Error,
			"error in updating the product",
		)
	}

	c.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	product, err = mapper.Map[*models.Product](productDataModel)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping Product",
		)
	}

	return product, err
}
