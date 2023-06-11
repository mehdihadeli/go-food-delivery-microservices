package infrastructure

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"
)

var Module = fx.Module("infrastructurefx",
	// infrastructure setup --> should go infrastructure module
	core.Module,
	customEcho.Module,
	grpc.Module,
	gormPostgres.Module,
	mongodb.Module,
	otel.Module,
	redis.Module,
	rabbitmq.Module(func(builder configurations.RabbitMQConfigurationBuilder) {
		fmt.Print("Creating")
	}),
)
