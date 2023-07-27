package infrastructure

import (
	"github.com/go-playground/validator"
	"go.uber.org/fx"

	rabbitmq2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/configurations/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/params"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/elasticsearch"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
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
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigOrdersRabbitMQ(builder)
			}
		},
	),

	// Other provides
	fx.Provide(validator.New),
)
