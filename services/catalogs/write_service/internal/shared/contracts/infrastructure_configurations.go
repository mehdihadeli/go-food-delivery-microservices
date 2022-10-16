package contracts

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"go.opentelemetry.io/otel/metric"
)

type InfrastructureConfigurations struct {
	Log             logger.Logger
	Cfg             *config.Config
	Validator       *validator.Validate
	Pgx             *postgres.Pgx
	Gorm            *gormPostgres.Gorm
	Metrics         metric.Meter
	Esdb            *esdb.Client
	ElasticClient   *elasticsearch.Client
	GrpcClient      grpc.GrpcClient
	EventSerializer serializer.EventSerializer
}
