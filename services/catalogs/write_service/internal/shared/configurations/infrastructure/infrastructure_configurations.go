package infrastructure

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"go.opentelemetry.io/otel/metric"
)

type infrastructureConfigurations struct {
	log             logger.Logger
	cfg             *config.Config
	validator       *validator.Validate
	pgx             *postgres.Pgx
	gorm            *gormPostgres.Gorm
	metrics         metric.Meter
	esdb            *esdb.Client
	elasticClient   *elasticsearch.Client
	grpcClient      grpc.GrpcClient
	eventSerializer serializer.EventSerializer
}

func (i *infrastructureConfigurations) Log() logger.Logger {
	return i.log
}

func (i *infrastructureConfigurations) Cfg() *config.Config {
	return i.cfg
}

func (i *infrastructureConfigurations) Validator() *validator.Validate {
	return i.validator
}

func (i *infrastructureConfigurations) Pgx() *postgres.Pgx {
	return i.pgx
}

func (i *infrastructureConfigurations) Gorm() *gormPostgres.Gorm {
	return i.gorm
}

func (i *infrastructureConfigurations) Esdb() *esdb.Client {
	return i.esdb
}

func (i *infrastructureConfigurations) GrpcClient() grpc.GrpcClient {
	return i.grpcClient
}

func (i *infrastructureConfigurations) ElasticClient() *elasticsearch.Client {
	return i.elasticClient
}

func (i *infrastructureConfigurations) EventSerializer() serializer.EventSerializer {
	return i.eventSerializer
}

func (i *infrastructureConfigurations) Metrics() metric.Meter {
	return i.metrics
}
