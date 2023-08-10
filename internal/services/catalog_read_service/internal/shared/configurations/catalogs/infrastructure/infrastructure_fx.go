package infrastructure

import (
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"
	"go.uber.org/fx"

	rabbitmq2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/configurations/rabbitmq"
)

// https://pmihaylov.com/shared-components-go-microservices/
var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	gormPostgres.Module,
	mongodb.Module,
	otel.Module,
	redis.Module,
	rabbitmq.ModuleFunc(
		func(v *validator.Validate, l logger.Logger, tracer tracing.AppTracer) configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder, l, v, tracer)
			}
		},
	),

	// Other provides
	fx.Provide(validator.New),
)
