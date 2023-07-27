package bus

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
)

type Bus interface {
	producer.Producer
	consumer.BusControl
	consumer.ConsumerConnector
}
