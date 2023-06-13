package catalogs

import (
	"go.uber.org/fx"

	appconfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/catalogs/infrastructure"
)

var Module = fx.Module(
	"catalogsfx",
	// Shared Modules
	appconfig.Module,
	infrastructure.Module,

	// Features Modules
	products.Module,

	// Other provides
	fx.Provide(NewCatalogsServiceConfigurator),
)
