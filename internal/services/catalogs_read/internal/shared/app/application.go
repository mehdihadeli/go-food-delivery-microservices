package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/catalogs"
)

type CatalogsReadApplication struct {
	*catalogs.CatalogsServiceConfigurator
}

func NewCatalogsReadApplication(
	providers []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment config.Environment,
) *CatalogsReadApplication {
	app := fxapp.NewApplication(providers, options, logger, environment)
	return &CatalogsReadApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}
