package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
)

type InfrastructureConfiguration struct {
	Log               logger.Logger
	Cfg               *config.Config
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Pgx               *postgres.Pgx
	Gorm              *gormPostgres.Gorm
	Metrics           *CatalogsServiceMetrics
	Esdb              *esdb.Client
	ElasticClient     *v7.Client
	CustomMiddlewares cutomMiddlewares.CustomMiddlewares
}

type InfrastructureConfigurator interface {
	ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func())
}

type infrastructureConfigurator struct {
	log logger.Logger
	cfg *config.Config
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config) InfrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func()) {
	infrastructure := &InfrastructureConfiguration{Cfg: ic.cfg, Log: ic.log, Validator: validator.New()}

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

	var cleanup []func()

	gorm, err := ic.configGorm()
	if err != nil {
		return nil, err, nil
	}
	infrastructure.Gorm = gorm

	err, jaegerCleanup := ic.configJaeger()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, jaegerCleanup)

	pgx, err, postgresCleanup := ic.configPostgres()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, postgresCleanup)
	infrastructure.Pgx = pgx

	kafkaConn, kafkaProducer, err, kafkaCleanup := ic.configKafka(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, kafkaCleanup)
	infrastructure.KafkaConn = kafkaConn
	infrastructure.KafkaProducer = kafkaProducer

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			c()
		}
	}
}
