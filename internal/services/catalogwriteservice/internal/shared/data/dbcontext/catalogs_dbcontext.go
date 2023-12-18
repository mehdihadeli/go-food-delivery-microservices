package dbcontext

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

type CatalogContextActionFunc func(ctx context.Context, catalogContext *CatalogsGormDBContext) error

type CatalogsGormDBContext struct {
	// our dbcontext base
	contracts.IGormDBContext
}

func NewCatalogsDBContext(db *gorm.DB) *CatalogsGormDBContext {
	// initialize base
	c := &CatalogsGormDBContext{IGormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
