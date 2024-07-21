package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/bus"
	messageConsumer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
	rabbitmq2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	rabbitmqConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/messaging/consumer"
	testUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_RabbitMQ_Container(t *testing.T) {
	ctx := context.Background()
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()

	var rabbitmqbus bus.Bus

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		rabbitmq2.ModuleFunc(
			func(l logger.Logger) rabbitmqConfigurations.RabbitMQConfigurationBuilderFuc {
				return func(builder configurations.RabbitMQConfigurationBuilder) {
					builder.AddConsumer(
						ProducerConsumerMessage{},
						func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
							consumerBuilder.WithHandlers(
								func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
									handlerBuilder.AddHandler(fakeConsumer)
								},
							)
						},
					)
				}
			},
		),
		fx.Decorate(RabbitmqDockerTestContainerOptionsDecorator(t, ctx)),
		fx.Populate(&rabbitmqbus),
	).RequireStart()

	require.NotNil(t, rabbitmqbus)

	err := rabbitmqbus.Start(ctx)
	require.NoError(t, err)

	//// wait for consumers ready to consume before publishing messages (for preventing messages lost)
	time.Sleep(time.Second * 2)

	err = rabbitmqbus.PublishMessage(
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
