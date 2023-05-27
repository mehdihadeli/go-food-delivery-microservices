package contracts

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
)

type InfrastructureConfigurations struct {
	Log             logger.Logger
	Cfg             *config.Config
	Validator       *validator.Validate
	Pgx             *postgres.Pgx
	Gorm            *gorm.DB
	Metrics         metric.Meter
	Esdb            *esdb.Client
	ElasticClient   *elasticsearch.Client
	GrpcClient      grpc.GrpcClient
	EventSerializer serializer.EventSerializer
}
