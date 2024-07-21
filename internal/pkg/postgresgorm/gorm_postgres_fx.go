package postgresgorm

import (
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/health/contracts"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"gormpostgresfx",
	fx.Provide(
		provideConfig,
		NewGorm,
		NewSQLDB,

		fx.Annotate(
			NewGormHealthChecker,
			fx.As(new(contracts.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	),
)
