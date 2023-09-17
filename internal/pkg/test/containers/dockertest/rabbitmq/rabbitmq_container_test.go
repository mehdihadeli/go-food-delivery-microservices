package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	messageConsumer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	rabbitmqConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_RabbitMQ_Container(t *testing.T) {
	ctx := context.Background()
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()
	defaultLogger.SetupDefaultLogger()
	eventSerializer := serializer.NewDefaultEventSerializer(json.NewDefaultSerializer())

	rabbitmq, err := NewRabbitMQDockerTest(
		defaultLogger.Logger,
	).Start(ctx, t, eventSerializer, func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
		builder.AddConsumer(ProducerConsumerMessage{},
			func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
				consumerBuilder.WithHandlers(
					func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
						handlerBuilder.AddHandler(fakeConsumer)
					},
				)
			})
	})

	require.NoError(t, err)
	require.NotNil(t, rabbitmq)

	err = rabbitmq.Start(ctx)
	require.NoError(t, err)

	// wait for consumers ready to consume before publishing messages (for preventing messages lost)
	time.Sleep(time.Second * 1)

	err = rabbitmq.PublishMessage(
		context.Background(),
		&ProducerConsumerMessage{
			Data:    "ssssssssss",
			Message: types.NewMessage(uuid.NewV4().String()),
		},
		nil,
	)
	if err != nil {
		return
	}

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	})

	t.Log("stopping test container")

	if err != nil {
		require.FailNow(t, err.Error())
	}
}

type ProducerConsumerMessage struct {
	*types.Message
	Data string
}
