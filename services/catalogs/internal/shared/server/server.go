package server

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Server struct {
	Log               logger.Logger
	Cfg               *config.Config
	HttpServer        *http.Server
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Im                interceptors.InterceptorManager
	PgConn            *pgxpool.Pool
	Metrics           *shared.CatalogsServiceMetrics
	Echo              *echo.Echo
	MongoClient       *mongo.Client
	ElasticClient     *v7.Client
	Ctx               context.Context
	MiddlewareManager middlewares.MiddlewareManager
	DoneCh            chan struct{}
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{Log: log, Cfg: cfg}
}
