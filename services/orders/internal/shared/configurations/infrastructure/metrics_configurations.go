package infrastructure

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type OrdersServiceMetrics struct {
	SuccessGrpcRequests prometheus.Counter
	ErrorGrpcRequests   prometheus.Counter

	CreateOrderGrpcRequests prometheus.Counter

	CreateOrderHttpRequests prometheus.Counter

	SuccessHttpRequests prometheus.Counter
	ErrorHttpRequests   prometheus.Counter

	SuccessKafkaMessages prometheus.Counter
	ErrorKafkaMessages   prometheus.Counter

	CreateOrderKafkaMessages prometheus.Counter
}

func (ic *infrastructureConfigurator) configCatalogsMetrics() *OrdersServiceMetrics {
	cfg := ic.cfg
	return &OrdersServiceMetrics{
		SuccessGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of success grpc requests",
		}),
		ErrorGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of error grpc requests",
		}),
		CreateOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_order_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of create order grpc requests",
		}),
		CreateOrderKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_order_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create order kafka messages",
		}),
		SuccessKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of success kafka processed messages",
		}),
		ErrorKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of error kafka processed messages",
		}),
		CreateOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of create order http requests",
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
