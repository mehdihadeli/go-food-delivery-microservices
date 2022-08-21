package data

import "context"

type GenericRepository[D DataModel[E], E any] interface {
	Add(ctx context.Context, entity E) error
	AddAll(ctx context.Context, entities []E) error
	GetById(ctx context.Context, id int) (E, error)
	GetAll(ctx context.Context) ([]E, error)
	Where(ctx context.Context, params E) ([]E, error)
	Update(ctx context.Context, entity E) error
	UpdateAll(ctx context.Context, entities []E) error
	Delete(ctx context.Context, id int) error
	SkipTake(skip int, take int, ctx context.Context) ([]E, error)
	Count(ctx context.Context) int64
	Find(ctx context.Context, specification Specification) ([]E, error)
}
