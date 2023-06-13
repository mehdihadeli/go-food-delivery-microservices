package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs"
)

type CatalogsWriteApplication struct {
	*catalogs.CatalogsConfigurator
}

func NewCatalogsWriteApplication(
	providers []interface{},
	options []fx.Option,
) *CatalogsWriteApplication {
	app := fxapp.NewApplication(providers, options)
	return &CatalogsWriteApplication{
		CatalogsConfigurator: catalogs.NewCatalogsConfigurator(app),
	}
}
