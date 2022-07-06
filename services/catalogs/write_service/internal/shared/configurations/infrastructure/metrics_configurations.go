package infrastructure

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CatalogsServiceMetrics struct {
	SuccessGrpcRequests prometheus.Counter
	ErrorGrpcRequests   prometheus.Counter

	CreateProductGrpcRequests  prometheus.Counter
	UpdateProductGrpcRequests  prometheus.Counter
	DeleteProductGrpcRequests  prometheus.Counter
	GetProductByIdGrpcRequests prometheus.Counter
	SearchProductGrpcRequests  prometheus.Counter

	CreateProductHttpRequests  prometheus.Counter
	UpdateProductHttpRequests  prometheus.Counter
	DeleteProductHttpRequests  prometheus.Counter
	GetProductByIdHttpRequests prometheus.Counter
	GetProductsHttpRequests    prometheus.Counter
	SearchProductHttpRequests  prometheus.Counter

	SuccessHttpRequests prometheus.Counter
	ErrorHttpRequests   prometheus.Counter

	SuccessKafkaMessages prometheus.Counter
	ErrorKafkaMessages   prometheus.Counter

	CreateProductKafkaMessages prometheus.Counter
	UpdateProductKafkaMessages prometheus.Counter
	DeleteProductKafkaMessages prometheus.Counter
}

func (ic *infrastructureConfigurator) configCatalogsMetrics() *CatalogsServiceMetrics {
	cfg := ic.cfg
	return &CatalogsServiceMetrics{
		SuccessGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of success grpc requests",
		}),
		ErrorGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of error grpc requests",
		}),
		CreateProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of create product grpc requests",
		}),
		UpdateProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update product grpc requests",
		}),
		DeleteProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of delete product grpc requests",
		}),
		GetProductByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_product_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get product by id grpc requests",
		}),
		SearchProductGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_product_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of search product grpc requests",
		}),
		CreateProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create product kafka messages",
		}),
		UpdateProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update product kafka messages",
		}),
		DeleteProductKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of delete product kafka messages",
		}),
		SuccessKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of success kafka processed messages",
		}),
		ErrorKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of error kafka processed messages",
		}),
		CreateProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of create product http requests",
		}),
		UpdateProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of update product http requests",
		}),
		DeleteProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete product http requests",
		}),
		GetProductByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_product_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get product by id http requests",
		}),
		SearchProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of search product http requests",
		}),
		GetProductsHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_products_http_requests_total", cfg.ServiceName),
			Help: "The total number of get products http requests",
		}),
		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
	}
}
