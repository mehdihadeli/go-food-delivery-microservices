package rabbitmq

import (
	"context"
	messageConsumer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_RabbitMQ_Container(t *testing.T) {
	ctx := context.Background()
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler()

	rabbitmq, err := NewRabbitMQTestContainers().Start(ctx, t, func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
		builder.AddConsumer(ProducerConsumerMessage{},
			func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
				consumerBuilder.WithHandlers(func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
					handlerBuilder.AddHandler(fakeConsumer)
				})
			})
	})

	require.NoError(t, err)
	require.NotNil(t, rabbitmq)

	err = rabbitmq.Start(ctx)
	require.NoError(t, err)

	// wait for consumers ready to consume before publishing messages (for preventing messages lost)
	time.Sleep(time.Second * 1)

	err = rabbitmq.PublishMessage(context.Background(), &ProducerConsumerMessage{Data: "ssssssssss", Message: types2.NewMessage(uuid.NewV4().String())}, nil)
	if err != nil {
		return
	}

	err = test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	})

	t.Log("stopping test container")

	if err != nil {
		require.FailNow(t, err.Error())
	}
}

type ProducerConsumerMessage struct {
	*types2.Message
	Data string
}
