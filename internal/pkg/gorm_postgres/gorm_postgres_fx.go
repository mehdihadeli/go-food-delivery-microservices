package gormPostgres

import (
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"

	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
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
	fx.Invoke(EnableTracing),
)

func EnableTracing(gormDb *gorm.DB) error {
	// add tracing to gorm
	return gormDb.Use(tracing.NewPlugin())
}
