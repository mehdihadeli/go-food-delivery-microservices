package test

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/test"
	constants2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/constants"

	"github.com/spf13/viper"
	"go.uber.org/fx/fxtest"
)

type CatalogsWriteTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

func NewCatalogsWriteTestApplicationBuilder(tb fxtest.TB) *CatalogsWriteTestApplicationBuilder {
	// set viper internal registry, in real app we read it from `.env` file in current `executing working directory` for example `catalogs_service`
	viper.Set(constants.PROJECT_NAME_ENV, constants2.PROJECT_NAME)

	return &CatalogsWriteTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

func (a *CatalogsWriteTestApplicationBuilder) Build() *CatalogsWriteTestApplication {
	return NewCatalogsWriteTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
