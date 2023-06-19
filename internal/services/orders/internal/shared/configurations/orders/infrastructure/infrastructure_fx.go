package infrastructure

import (
	"github.com/go-playground/validator"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/elasticsearch"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	rabbitmq2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts/params"
)

// https://pmihaylov.com/shared-components-go-microservices/

var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	mongodb.Module,
	elasticsearch.Module,
	eventstroredb.Module(
		func(params params.OrderProjectionParams) eventstroredb.ProjectionBuilderFuc {
			return func(builder eventstroredb.ProjectionsBuilder) {
				builder.AddProjections(params.Projections)
			}
		},
	),
	otel.Module,
	rabbitmq.Module(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigOrdersRabbitMQ(builder)
			}
		},
	),

	// Other provides
	fx.Provide(validator.New),
)
