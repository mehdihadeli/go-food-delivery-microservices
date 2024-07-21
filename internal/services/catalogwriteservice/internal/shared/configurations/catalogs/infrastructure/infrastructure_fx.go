package infrastructure

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/health"
	customEcho "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration/goose"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresmessaging"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	rabbitmq2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations/rabbitmq"

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
	postgresgorm.Module,
	postgresmessaging.Module,
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
