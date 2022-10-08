package contracts

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"go.opentelemetry.io/otel/metric"
)

type InfrastructureConfigurations interface {
	Log() logger.Logger
	Cfg() *config.Config
	Validator() *validator.Validate
	Pgx() *postgres.Pgx
	Gorm() *gormPostgres.Gorm
	Esdb() *esdb.Client
	Metrics() metric.Meter
	GrpcClient() grpc.GrpcClient
	EventSerializer() serializer.EventSerializer
}
