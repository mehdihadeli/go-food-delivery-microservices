package data

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"datafx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		dbcontext.NewCatalogsDBContext,
	),
)
