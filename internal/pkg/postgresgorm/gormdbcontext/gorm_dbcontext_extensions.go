package gormdbcontext

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	defaultlogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/scopes"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
	uuid "github.com/satori/go.uuid"
)

func Exists[TDataModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	id uuid.UUID,
) bool {
	var count int64

	dataModel := typeMapper.GenericInstanceByT[TDataModel]()

	dbContext.DB().WithContext(ctx).Model(dataModel).Scopes(scopes.FilterByID(id)).Count(&count)

	return count > 0
}

func FindModelByID[TDataModel interface{}, TModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	id uuid.UUID,
) (TModel, error) {
	var dataModel TDataModel

	// https://gorm.io/docs/query.html#Retrieving-objects-with-primary-key
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	// https://gorm.io/docs/query.html#Inline-Condition
	// https://gorm.io/docs/advanced_query.html
	// result := c.WithContext(ctx).First(&dataModel, "id = ?", id)
	// result := c.WithContext(ctx).First(&TDataModel{Id: id})
	// result := c.WithContext(ctx).Scopes(scopes.FilterByID(id)).First(&dataModel)

	modelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TModel]())
	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())

	result := dbContext.DB().WithContext(ctx).First(&dataModel, id)
	if result.Error != nil {
		return *new(TModel), customErrors.NewNotFoundErrorWrap(
			result.Error,
			fmt.Sprintf(
				"%s with id `%s` not found in the database",
				dataModelName,
				id.String(),
			),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	resultModel, err := mapper.Map[TModel](dataModel)
	if err != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			err,
			fmt.Sprintf("error in the mapping %s", modelName),
		)
	}

	return resultModel, nil
}

func FindDataModelByID[TDataModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	id uuid.UUID,
) (TDataModel, error) {
	var dataModel TDataModel

	// https://gorm.io/docs/query.html#Retrieving-objects-with-primary-key
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	// https://gorm.io/docs/query.html#Inline-Condition
	// https://gorm.io/docs/advanced_query.html
	// result := c.WithContext(ctx).First(&dataModel, "id = ?", id)
	// result := c.WithContext(ctx).First(&TDataModel{Id: id})
	// result := c.WithContext(ctx).Scopes(scopes.FilterByID(id)).First(&dataModel)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())

	result := dbContext.DB().WithContext(ctx).First(&dataModel, id)
	if result.Error != nil {
		return *new(TDataModel), customErrors.NewNotFoundErrorWrap(
			result.Error,
			fmt.Sprintf(
				"%s with id `%s` not found in the database",
				dataModelName,
				id.String(),
			),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	return dataModel, nil
}

// DeleteDataModelByID delete the data-model inner a tx if exists
func DeleteDataModelByID[TDataModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	id uuid.UUID,
) error {
	txDBContext := dbContext.WithTxIfExists(ctx)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())

	exists := Exists[TDataModel](ctx, dbContext, id)
	if !exists {
		return customErrors.NewNotFoundError(fmt.Sprintf("%s with id `%s` not found in the database",
			dataModelName,
			id.String(),
		))
	}

	dataModel := typeMapper.GenericInstanceByT[TDataModel]()

	// https://gorm.io/docs/delete.html#Delete-a-Record
	// https://gorm.io/docs/delete.html#Find-soft-deleted-records
	// result := dbContext.WithContext(ctx).Delete(&TDataModel{Id: id})
	result := txDBContext.DB().WithContext(ctx).Delete(dataModel, id)
	if result.Error != nil {
		return customErrors.NewInternalServerErrorWrap(
			result.Error,
			fmt.Sprintf(
				"error in deleting %s with id `%s` in the database",
				dataModelName,
				id.String(),
			),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	return nil
}

// AddModel add the model inner a tx if exists
func AddModel[TDataModel interface{}, TModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	model TModel,
) (TModel, error) {
	txDBContext := dbContext.WithTxIfExists(ctx)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())
	modelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TModel]())

	dataModel, err := mapper.Map[TDataModel](model)
	if err != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			err,
			fmt.Sprintf("error in the mapping %s", dataModelName),
		)
	}

	// https://gorm.io/docs/create.html
	result := txDBContext.DB().WithContext(ctx).Create(dataModel)
	if result.Error != nil {
		return *new(TModel), customErrors.NewConflictErrorWrap(
			result.Error,
			fmt.Sprintf("%s already exists", modelName),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	resultModel, err := mapper.Map[TModel](dataModel)
	if err != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			err,
			fmt.Sprintf("error in the mapping %s", modelName),
		)
	}

	return resultModel, err
}

// AddDataModel add the data-model inner a tx if exists
func AddDataModel[TDataModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	dataModel TDataModel,
) (TDataModel, error) {
	txDBContext := dbContext.WithTxIfExists(ctx)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())

	// https://gorm.io/docs/create.html
	result := txDBContext.DB().WithContext(ctx).Create(dataModel)
	if result.Error != nil {
		return *new(TDataModel), customErrors.NewConflictErrorWrap(
			result.Error,
			fmt.Sprintf("%s already exists", dataModelName),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	return dataModel, nil
}

// UpdateModel update the model inner a tx if exists
func UpdateModel[TDataModel interface{}, TModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	model TModel,
) (TModel, error) {
	txDBContext := dbContext.WithTxIfExists(ctx)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())
	modelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TModel]())

	dataModel, err := mapper.Map[TDataModel](model)
	if err != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			err,
			fmt.Sprintf("error in the mapping %s", dataModelName),
		)
	}

	// https://gorm.io/docs/update.html
	result := txDBContext.DB().WithContext(ctx).Updates(dataModel)
	if result.Error != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			result.Error,
			fmt.Sprintf("error in updating the %s", modelName),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	modelResult, err := mapper.Map[TModel](dataModel)
	if err != nil {
		return *new(TModel), customErrors.NewInternalServerErrorWrap(
			err,
			fmt.Sprintf("error in the mapping %s", modelName),
		)
	}

	return modelResult, err
}

// UpdateDataModel update the data-model inner a tx if exists
func UpdateDataModel[TDataModel interface{}](
	ctx context.Context,
	dbContext contracts.GormDBContext,
	dataModel TDataModel,
) (TDataModel, error) {
	txDBContext := dbContext.WithTxIfExists(ctx)

	dataModelName := strcase.ToSnake(typeMapper.GetGenericNonePointerTypeNameByT[TDataModel]())

	// https://gorm.io/docs/update.html
	result := txDBContext.DB().WithContext(ctx).Updates(dataModel)
	if result.Error != nil {
		return *new(TDataModel), customErrors.NewInternalServerErrorWrap(
			result.Error,
			fmt.Sprintf("error in updating the %s", dataModelName),
		)
	}

	defaultlogger.GetLogger().Infof("Number of affected rows are: %d", result.RowsAffected)

	return dataModel, nil
}
