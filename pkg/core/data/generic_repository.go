package data

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type GenericRepositoryWithDataModel[TDataModel interface{}, TEntity interface{}] interface {
	Add(ctx context.Context, entity TEntity) error
	AddAll(ctx context.Context, entities []TEntity) error
	GetById(ctx context.Context, id uuid.UUID) (TEntity, error)
	GetAll(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[TEntity], error)
	Search(ctx context.Context, searchTerm string, listQuery *utils.ListQuery) (*utils.ListResult[TEntity], error)
	Where(ctx context.Context, filters map[string]interface{}) ([]TEntity, error)
	Update(ctx context.Context, entity TEntity) error
	UpdateAll(ctx context.Context, entities []TEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	SkipTake(skip int, take int, ctx context.Context) ([]TEntity, error)
	Count(ctx context.Context) int64
	Find(ctx context.Context, specification Specification) ([]TEntity, error)
}

type GenericRepository[TEntity interface{}] interface {
	GenericRepositoryWithDataModel[TEntity, TEntity]
}
