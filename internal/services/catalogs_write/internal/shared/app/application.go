package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs"
)

type CatalogsWriteApplication struct {
	*catalogs.CatalogsServiceConfigurator
}

func NewCatalogsWriteApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environemnt.Environment,
) *CatalogsWriteApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &CatalogsWriteApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}
