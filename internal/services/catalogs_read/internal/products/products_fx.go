package products

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/data/repositories"
)

var Module = fx.Module(
	"productsfx",

	// Other provides
	fx.Provide(repositories.NewRedisProductRepository),
	fx.Provide(repositories.NewMongoProductRepository),
)
