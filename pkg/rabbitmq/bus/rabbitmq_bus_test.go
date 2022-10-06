package bus

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer/json"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	messageConsumer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/pipeline"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AddRabbitMQ(t *testing.T) {
	ctx := context.Background()

	fakeConsumer := consumer.NewRabbitMQFakeTestConsumer()
	b, err := NewRabbitMQBus(ctx, &config.RabbitMQConfig{
		RabbitMqHostOptions: &config.RabbitMqHostOptions{
			UserName: "guest",
			Password: "guest",
			HostName: "localhost",
			Port:     5672,
		}},
		func(builder configurations.RabbitMQConfigurationBuilder) {
			builder.AddProducer(ProducerConsumerMessage{}, func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
			})
			builder.AddConsumer(ProducerConsumerMessage{},
				func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
					builder.WithHandlers(func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
						consumerHandlerBuilder.AddHandler(NewTestMessageHandler())
						consumerHandlerBuilder.AddHandler(NewTestMessageHandler2())
						consumerHandlerBuilder.AddHandler(fakeConsumer)
					}).WIthPipelines(func(consumerPipelineBuilder pipeline.ConsumerPipelineConfigurationBuilder) {
						consumerPipelineBuilder.AddPipeline(NewPipeline1())
					})
				})
		}, json.NewJsonEventSerializer(), defaultLogger.Logger)
	if err != nil {
		return
	}

	err = b.Start(ctx)
	if err != nil {
		return
	}

	err = b.PublishMessage(context.Background(), &ProducerConsumerMessage{Data: "ssssssssss", Message: types2.NewMessage(uuid.NewV4().String())}, nil)
	if err != nil {
		return
	}

	err = test.WaitUntilConditionMet(func() bool {
		return fakeConsumer.IsHandled()
	})
	assert.NoError(t, err)

	b.Stop(ctx)
}

type ProducerConsumerMessage struct {
	*types2.Message
	Data string
}

func NewProducerConsumerMessage(data string) *ProducerConsumerMessage {
	return &ProducerConsumerMessage{
		Data:    data,
		Message: types2.NewMessage(uuid.NewV4().String()),
	}
}

// /////////// ConsumerHandlerT
type TestMessageHandler struct {
}

func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

func (t *TestMessageHandler) Handle(ctx context.Context, consumeContext types2.MessageConsumeContext) error {
	message := consumeContext.Message().(*ProducerConsumerMessage)
	fmt.Println(message)

	return nil
}

type TestMessageHandler2 struct {
}

func (t *TestMessageHandler2) Handle(ctx context.Context, consumeContext types2.MessageConsumeContext) error {
	message := consumeContext.Message()
	fmt.Println(message)

	return nil
}

func NewTestMessageHandler2() *TestMessageHandler2 {
	return &TestMessageHandler2{}
}

// /////////////// ConsumerPipeline
type Pipeline1 struct {
}

func NewPipeline1() pipeline.ConsumerPipeline {
	return &Pipeline1{}
}

func (p Pipeline1) Handle(ctx context.Context, consumerContext types2.MessageConsumeContext, next pipeline.ConsumerHandlerFunc) error {
	fmt.Println("PipelineBehaviourTest.Handled")

	fmt.Println(fmt.Sprintf("pipeline got a message with id '%s'", consumerContext.Message().GeMessageId()))

	err := next()
	if err != nil {
		return err
	}
	return nil
}
