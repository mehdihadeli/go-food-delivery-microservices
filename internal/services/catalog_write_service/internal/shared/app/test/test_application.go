package test

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/test"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
)

type CatalogsWriteTestApplication struct {
	*app.CatalogsWriteApplication
	tb fxtest.TB
}

func NewCatalogsWriteTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environemnt.Environment,
) *CatalogsWriteTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		logger,
		environment,
	)

	catalogApplication := &app.CatalogsWriteApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(testApp),
	}

	return &CatalogsWriteTestApplication{
		CatalogsWriteApplication: catalogApplication,
		tb:                       tb,
	}
}
