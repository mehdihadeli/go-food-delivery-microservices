package bus

import (
	"context"
	"fmt"
	"testing"

	messageConsumer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	pipeline2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/pipeline"
	types3 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer/json"
	defaultlogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	rabbitmqconsumer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer"
	consumerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	rabbitmqproducer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer"
	producerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddRabbitMQ(t *testing.T) {
	testUtils.SkipCI(t)
	ctx := context.Background()

	fakeConsumer2 := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()
	fakeConsumer3 := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()

	serializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)

	//rabbitmqOptions := &config.RabbitmqOptions{
	//	RabbitmqHostOptions: &config.RabbitmqHostOptions{
	//		UserName: "guest",
	//		Password: "guest",
	//		HostName: "localhost",
	//		Port:     5672,
	//	},
	//}

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(defaultlogger.GetLogger()).
		PopulateContainerOptions(ctx, t)
	require.NoError(t, err)

	options := &config.RabbitmqOptions{
		RabbitmqHostOptions: rabbitmqHostOption,
	}

	conn, err := types.NewRabbitMQConnection(options)
	require.NoError(t, err)

	consumerFactory := rabbitmqconsumer.NewConsumerFactory(
		options,
		conn,
		serializer,
		defaultlogger.GetLogger(),
	)
	producerFactory := rabbitmqproducer.NewProducerFactory(
		options,
		conn,
		serializer,
		defaultlogger.GetLogger(),
	)

	b, err := NewRabbitmqBus(
		defaultlogger.GetLogger(),
		consumerFactory,
		producerFactory,
		func(builder configurations.RabbitMQConfigurationBuilder) {
			builder.AddProducer(
				ProducerConsumerMessage{},
				func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
				},
			)
			builder.AddConsumer(
				ProducerConsumerMessage{},
				func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
					builder.WithHandlers(func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
						consumerHandlerBuilder.AddHandler(
							NewTestMessageHandler(),
						)
						consumerHandlerBuilder.AddHandler(
							NewTestMessageHandler2(),
						)
					}).
						WIthPipelines(func(consumerPipelineBuilder pipeline2.ConsumerPipelineConfigurationBuilder) {
							consumerPipelineBuilder.AddPipeline(NewPipeline1())
						})
				},
			)
		},
	)

	require.NoError(t, err)

	//err = b.ConnectRabbitMQConsumer(ProducerConsumerMessage{}, func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
	//	consumerBuilder.WithHandlers(func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
	//		handlerBuilder.AddHandler(fakeConsumer)
	//	})
	//})
	//require.NoError(t, err)

	err = b.ConnectConsumerHandler(&ProducerConsumerMessage{}, fakeConsumer2)
	require.NoError(t, err)

	err = b.ConnectConsumerHandler(&ProducerConsumerMessage{}, fakeConsumer3)
	require.NoError(t, err)

	err = b.Start(ctx)
	require.NoError(t, err)

	err = b.PublishMessage(
		context.Background(),
		&ProducerConsumerMessage{
			Data:    "ssssssssss",
			Message: types3.NewMessage(uuid.NewV4().String()),
		},
		nil,
	)
	require.NoError(t, err)

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer2.IsHandled() && fakeConsumer3.IsHandled()
	})
	assert.NoError(t, err)

	b.Stop()
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

func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

func (t *TestMessageHandler) Handle(
	ctx context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	message := consumeContext.Message().(*ProducerConsumerMessage)
	fmt.Println(message)

	return nil
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

func NewPipeline1() pipeline2.ConsumerPipeline {
	return &Pipeline1{}
}

func (p *Pipeline1) Handle(
	ctx context.Context,
	consumerContext types3.MessageConsumeContext,
	next pipeline2.ConsumerHandlerFunc,
) error {
	fmt.Println("PipelineBehaviourTest.Handled")

	fmt.Printf(
		"pipeline got a message with id '%s'",
		consumerContext.Message().GeMessageId(),
	)

	err := next(ctx)
	if err != nil {
		return err
	}

	return nil
}
