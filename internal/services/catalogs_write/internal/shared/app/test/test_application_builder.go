package test

import (
	"github.com/spf13/viper"
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/test"
	constants2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/constants"
)

type CatalogsWriteTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

func NewCatalogsWriteTestApplicationBuilder(tb fxtest.TB) *CatalogsWriteTestApplicationBuilder {
	// set viper internal registry
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
