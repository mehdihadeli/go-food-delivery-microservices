package application

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/catalogs"
)

type CatalogReadApplicationBuilder struct {
	*fxapp.ApplicationBuilder
}

func NewCatalogReadApplicationBuilder() *CatalogReadApplicationBuilder {
	return &CatalogReadApplicationBuilder{fxapp.NewApplicationBuilder()}
}

func (a *CatalogReadApplicationBuilder) Build() *CatalogReadApplication {
	app := fxapp.NewApplication(a.Providers, a.Options)
	return &CatalogReadApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}

type CatalogReadApplication struct {
	*catalogs.CatalogsServiceConfigurator
}
