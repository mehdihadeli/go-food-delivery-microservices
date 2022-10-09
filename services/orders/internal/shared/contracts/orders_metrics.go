package contracts

import "go.opentelemetry.io/otel/metric/instrument/syncfloat64"

type OrdersMetrics interface {
	SuccessGrpcRequests() syncfloat64.Counter
	ErrorGrpcRequests() syncfloat64.Counter

	CreateOrderGrpcRequests() syncfloat64.Counter
	UpdateOrderGrpcRequests() syncfloat64.Counter
	PayOrderGrpcRequests() syncfloat64.Counter
	SubmitOrderGrpcRequests() syncfloat64.Counter
	GetOrderByIdGrpcRequests() syncfloat64.Counter
	GetOrdersGrpcRequests() syncfloat64.Counter
	SearchOrderGrpcRequests() syncfloat64.Counter

	SuccessHttpRequests() syncfloat64.Counter
	ErrorHttpRequests() syncfloat64.Counter

	CreateOrderHttpRequests() syncfloat64.Counter
	UpdateOrderHttpRequests() syncfloat64.Counter
	PayOrderHttpRequests() syncfloat64.Counter
	SubmitOrderHttpRequests() syncfloat64.Counter
	GetOrderByIdHttpRequests() syncfloat64.Counter
	SearchOrderHttpRequests() syncfloat64.Counter
	GetOrdersHttpRequests() syncfloat64.Counter

	SuccessRabbitMQMessages() syncfloat64.Counter
	ErrorRabbitMQMessages() syncfloat64.Counter

	CreateOrderRabbitMQMessages() syncfloat64.Counter
	UpdateOrderRabbitMQMessages() syncfloat64.Counter
	DeleteOrderRabbitMQMessages() syncfloat64.Counter
}
