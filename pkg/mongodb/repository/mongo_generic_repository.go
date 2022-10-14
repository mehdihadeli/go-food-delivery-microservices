package repository

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGenericRepository[TDataModel interface{}, TEntity interface{}] struct {
	db *mongo.Client
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Add(ctx context.Context, entity TEntity) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) AddAll(ctx context.Context, entities []TEntity) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) GetById(ctx context.Context, id uuid.UUID) (TEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) GetAll(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[TEntity], error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Search(ctx context.Context, searchTerm string, listQuery *utils.ListQuery) (*utils.ListResult[TEntity], error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Where(ctx context.Context, filters map[string]interface{}) ([]TEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Update(ctx context.Context, entity TEntity) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) UpdateAll(ctx context.Context, entities []TEntity) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) SkipTake(skip int, take int, ctx context.Context) ([]TEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Count(ctx context.Context) int64 {
	//TODO implement me
	panic("implement me")
}

func (m *mongoGenericRepository[TDataModel, TEntity]) Find(ctx context.Context, specification data.Specification) ([]TEntity, error) {
	//TODO implement me
	panic("implement me")
}
