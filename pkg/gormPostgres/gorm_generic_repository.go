package gormPostgres

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	reflectionHelper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/reflection_helper"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"gorm.io/gorm"
)

// gorm generic repository
type gormGenericRepository[D data.DataModel[E], E any] struct {
	db *gorm.DB
}

// NewGenericGormRepository create new gorm generic repository
func NewGenericGormRepository[D data.DataModel[E], E any](db *gorm.DB) *gormGenericRepository[D, E] {
	return &gormGenericRepository[D, E]{
		db: db,
	}
}

func (r *gormGenericRepository[D, E]) Add(ctx context.Context, entity E) error {
	var dataModel D
	typeName := typeMapper.GetFullTypeName(dataModel)
	dataModel = typeMapper.GenericInstanceByTypeName[D](typeName)
	dataModel.FromEntity(entity)

	err := r.db.WithContext(ctx).Create(&dataModel).Error
	if err != nil {
		return err
	}

	reflectionHelper.SetValue[E](entity, dataModel.ToEntity())

	return nil
}

func (r *gormGenericRepository[D, E]) AddAll(ctx context.Context, entities []E) error {
	for _, entity := range entities {
		err := r.Add(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *gormGenericRepository[D, E]) GetById(ctx context.Context, id int) (E, error) {
	var dataModel D
	err := r.db.WithContext(ctx).First(&dataModel, id).Error
	if err != nil {
		return *new(E), err
	}

	return dataModel.ToEntity(), nil
}

func (r *gormGenericRepository[D, E]) GetAll(ctx context.Context) ([]E, error) {
	var dataModels []D
	err := r.db.WithContext(ctx).Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	var entities []E
	for _, dataModel := range dataModels {
		entities = append(entities, dataModel.ToEntity())
	}

	return entities, nil
}

func (r *gormGenericRepository[D, E]) Where(ctx context.Context, params E) ([]E, error) {
	var dataModels []D
	err := r.db.WithContext(ctx).Where(&params).Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	var entities []E
	for _, dataModel := range dataModels {
		entities = append(entities, dataModel.ToEntity())
	}

	return entities, nil
}

func (r *gormGenericRepository[D, E]) Update(ctx context.Context, entity E) error {
	var dataModel D
	typeName := typeMapper.GetFullTypeName(dataModel)
	dataModel = typeMapper.GenericInstanceByTypeName[D](typeName)
	dataModel.FromEntity(entity)

	err := r.db.WithContext(ctx).Save(&dataModel).Error
	if err != nil {
		return err
	}
	reflectionHelper.SetValue[E](entity, dataModel.ToEntity())

	return nil
}

func (r gormGenericRepository[D, E]) UpdateAll(ctx context.Context, entities []E) error {
	for _, e := range entities {
		err := r.Update(ctx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *gormGenericRepository[D, E]) Delete(ctx context.Context, id int) error {
	var dataModel D
	err := r.db.WithContext(ctx).First(&dataModel, id).Error
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Delete(&dataModel).Error
}

func (r *gormGenericRepository[D, E]) SkipTake(skip int, take int, ctx context.Context) ([]E, error) {
	var dataModels []D
	err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	var entities []E
	for _, dataModel := range dataModels {
		entities = append(entities, dataModel.ToEntity())
	}

	return entities, nil
}

func (r *gormGenericRepository[D, E]) Count(ctx context.Context) int64 {
	var dataModel D
	var count int64
	r.db.WithContext(ctx).Model(&dataModel).Count(&count)
	return count
}

func (r *gormGenericRepository[M, E]) Find(ctx context.Context, specification data.Specification) ([]E, error) {
	var models []M
	err := r.db.WithContext(ctx).Where(specification.GetQuery(), specification.GetValues()...).Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}
