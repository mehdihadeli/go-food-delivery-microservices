package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/catalogs"
)

type CatalogsReadApplication struct {
	*catalogs.CatalogsConfigurator
}

func NewCatalogsReadApplication(
	providers []interface{},
	options []fx.Option,
) *CatalogsReadApplication {
	app := fxapp.NewApplication(providers, options)
	return &CatalogsReadApplication{
		CatalogsConfigurator: catalogs.NewCatalogsConfigurator(app),
	}
}
