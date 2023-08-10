package infrastructure

import (
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	"go.uber.org/fx"

	rabbitmq2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/configurations/rabbitmq"
)

// https://pmihaylov.com/shared-components-go-microservices/

var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	gormPostgres.Module,
	postgres.Module,
	otel.Module,
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder)
			}
		},
	),

	// Other provides
	fx.Provide(validator.New),
)
