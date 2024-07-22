package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"

	messageConsumer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/pipeline"
	types3 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer/json"
	defaultLogger2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_Consumer_With_Fake_Message(t *testing.T) {
	testUtils.SkipCI(t)

	ctx := context.Background()

	//options := &config.RabbitmqOptions{
	//	RabbitmqHostOptions: &config.RabbitmqHostOptions{
	//		UserName: "guest",
	//		Password: "guest",
	//		HostName: "localhost",
	//		Port:     5672,
	//	},
	//}

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(defaultLogger2.GetLogger()).
		PopulateContainerOptions(ctx, t)
	require.NoError(t, err)

	options := &config.RabbitmqOptions{
		RabbitmqHostOptions: rabbitmqHostOption,
	}

	conn, err := types.NewRabbitMQConnection(options)
	require.NoError(t, err)

	eventSerializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)
	consumerFactory := NewConsumerFactory(
		options,
		conn,
		eventSerializer,
		defaultLogger2.GetLogger(),
	)
	producerFactory := producer.NewProducerFactory(
		options,
		conn,
		eventSerializer,
		defaultLogger2.GetLogger(),
	)

	fakeHandler := consumer.NewRabbitMQFakeTestConsumerHandler[ProducerConsumerMessage]()

	rabbitmqBus, err := bus.NewRabbitmqBus(
		defaultLogger2.GetLogger(),
		consumerFactory,
		producerFactory,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			builder.AddConsumer(
				ProducerConsumerMessage{},
				func(consumerBuilder configurations.RabbitMQConsumerConfigurationBuilder) {
					consumerBuilder.WithHandlers(
						func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
							consumerHandlerBuilder.AddHandler(fakeHandler)
						},
					)
				},
			)
		},
	)

	rabbitmqBus.Start(ctx)
	defer rabbitmqBus.Stop()

	time.Sleep(time.Second * 1)

	require.NoError(t, err)

	err = rabbitmqBus.PublishMessage(
		ctx,
		NewProducerConsumerMessage("test"),
		nil,
	)
	for err != nil {
		err = rabbitmqBus.PublishMessage(
			ctx,
			NewProducerConsumerMessage("test"),
			nil,
		)
	}

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeHandler.IsHandled()
	})

	require.NoError(t, err)
}

type ProducerConsumerMessage struct {
	*types3.Message
	Data string
}

func NewProducerConsumerMessage(data string) *ProducerConsumerMessage {
	return &ProducerConsumerMessage{
		Data:    data,
		Message: types3.NewMessage(uuid.NewV4().String()),
	}
}

// /////////// ConsumerHandlerT
type TestMessageHandler struct{}

func (t *TestMessageHandler) Handle(
	ctx context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	message := consumeContext.Message().(*ProducerConsumerMessage)
	fmt.Println(message)

	return nil
}

func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

type TestMessageHandler2 struct{}

func (t *TestMessageHandler2) Handle(
	ctx context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	message := consumeContext.Message()
	fmt.Println(message)

	return nil
}

func NewTestMessageHandler2() *TestMessageHandler2 {
	return &TestMessageHandler2{}
}

// /////////////// ConsumerPipeline
type Pipeline1 struct{}

func NewPipeline1() pipeline.ConsumerPipeline {
	return &Pipeline1{}
}

func (p Pipeline1) Handle(
	ctx context.Context,
	consumerContext types3.MessageConsumeContext,
	next pipeline.ConsumerHandlerFunc,
) error {
	fmt.Println("PipelineBehaviourTest.Handled")

	fmt.Println(
		fmt.Sprintf(
			"pipeline got a message with id '%s'",
			consumerContext.Message().GeMessageId(),
		),
	)

	err := next(ctx)
	if err != nil {
		return err
	}
	return nil
}
