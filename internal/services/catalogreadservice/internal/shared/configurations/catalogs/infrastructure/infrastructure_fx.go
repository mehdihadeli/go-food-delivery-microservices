package infrastructure

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/health"
	customEcho "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/redis"
	rabbitmq2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/configurations/rabbitmq"

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
	mongodb.Module,
	redis.Module,
	rabbitmq.ModuleFunc(
		func(v *validator.Validate, l logger.Logger, tracer tracing.AppTracer) configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder, l, v, tracer)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
