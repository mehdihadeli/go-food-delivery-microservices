package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/goose"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	rabbitmq2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/configurations/rabbitmq"

	"github.com/go-playground/validator"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/

var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	gormPostgres.Module,
	goose.Module,
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
