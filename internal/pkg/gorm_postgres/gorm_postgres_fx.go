package gormPostgres

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
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
			fx.As(new(health.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	),
)
