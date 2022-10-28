package bus

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
)

type Bus interface {
	producer.Producer
	consumer.BusControl
	consumer.ConsumerConnector
}
