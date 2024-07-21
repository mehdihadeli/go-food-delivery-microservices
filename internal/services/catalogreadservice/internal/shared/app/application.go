package app

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

	"go.uber.org/fx"
)

type CatalogsReadApplication struct {
	*catalogs.CatalogsServiceConfigurator
}

func NewCatalogsReadApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *CatalogsReadApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &CatalogsReadApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}
