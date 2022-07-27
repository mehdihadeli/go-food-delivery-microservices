package infrastructure

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type OrdersServiceMetrics struct {
	SuccessGrpcRequests prometheus.Counter
	ErrorGrpcRequests   prometheus.Counter

	CreateOrderGrpcRequests        prometheus.Counter
	UpdateOrderGrpcRequests        prometheus.Counter
	PayOrderGrpcRequests           prometheus.Counter
	SubmitOrderGrpcRequests        prometheus.Counter
	GetOrderByIdGrpcRequests       prometheus.Counter
	SearchOrderGrpcRequests        prometheus.Counter
	CancelOrderGrpcRequests        prometheus.Counter
	CompleteOrderGrpcRequests      prometheus.Counter
	ChangeAddressOrderGrpcRequests prometheus.Counter

	SuccessHttpRequests prometheus.Counter
	ErrorHttpRequests   prometheus.Counter

	CreateOrderHttpRequests        prometheus.Counter
	UpdateOrderHttpRequests        prometheus.Counter
	PayOrderHttpRequests           prometheus.Counter
	SubmitOrderHttpRequests        prometheus.Counter
	GetOrderByIdHttpRequests       prometheus.Counter
	SearchOrderHttpRequests        prometheus.Counter
	CompleteOrderHttpRequests      prometheus.Counter
	ChangeAddressOrderHttpRequests prometheus.Counter

	SuccessKafkaMessages prometheus.Counter
	ErrorKafkaMessages   prometheus.Counter

	CreateProductKafkaMessages prometheus.Counter
	UpdateProductKafkaMessages prometheus.Counter
	DeleteProductKafkaMessages prometheus.Counter
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
		UpdateOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_order_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update order grpc requests",
		}),
		PayOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_pay_order_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of pay order grpc requests",
		}),
		SubmitOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_submit_order_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of submit order grpc requests",
		}),
		GetOrderByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_order_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get order by id grpc requests",
		}),
		SearchOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_order_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of search order grpc requests",
		}),

		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
		CreateOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of create order http requests",
		}),
		UpdateOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of update order http requests",
		}),
		PayOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_pay_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of pay order http requests",
		}),
		SubmitOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_submit_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of submit order http requests",
		}),
		GetOrderByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_order_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get order by id http requests",
		}),
		SearchOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of search order http requests",
		}),
		CancelOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_cancel_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of cancel order http requests",
		}),
		CompleteOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_complete_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of complete order http requests",
		}),
		ChangeAddressOrderGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_change_address_order_gRPC_requests_total", cfg.ServiceName),
			Help: "The total number of change address order gRPC requests",
		}),
		ChangeAddressOrderHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_change_address_order_http_requests_total", cfg.ServiceName),
			Help: "The total number of change address order http requests",
		}),
	}
}
