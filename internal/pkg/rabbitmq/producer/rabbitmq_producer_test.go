package producer

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	types2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer/json"
	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	testUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_Publish_Message(t *testing.T) {
	testUtils.SkipCI(t)

	eventSerializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)

	ctx := context.Background()
	tp, err := tracing.NewOtelTracing(
		&tracing.TracingOptions{
			ServiceName:     "test",
			Enabled:         true,
			AlwaysOnSampler: true,
			ZipkinExporterOptions: &tracing.ZipkinExporterOptions{
				Url: "http://localhost:9411/api/v2/spans",
			},
		},
		environment.Development,
	)
	if err != nil {
		return
	}
	defer tp.Shutdown(ctx)

	//options := &config.RabbitmqOptions{
	//	RabbitmqHostOptions: &config.RabbitmqHostOptions{
	//		UserName: "guest",
	//		Password: "guest",
	//		HostName: "localhost",
	//		Port:     5672,
	//	},
	//}

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(ctx, t)
	require.NoError(t, err)

	options := &config.RabbitmqOptions{
		RabbitmqHostOptions: rabbitmqHostOption,
	}

	conn, err := types.NewRabbitMQConnection(options)
	require.NoError(t, err)

	producerFactory := NewProducerFactory(
		options,
		conn,
		eventSerializer,
		defaultLogger.GetLogger(),
	)

	rabbitmqProducer, err := producerFactory.CreateProducer(nil)

	require.NoError(t, err)

	err = rabbitmqProducer.PublishMessage(ctx, NewProducerMessage("test"), nil)
	require.NoError(t, err)
}

type ProducerMessage struct {
	*types2.Message
	Data string
}

func NewProducerMessage(data string) *ProducerMessage {
	return &ProducerMessage{
		Data:    data,
		Message: types2.NewMessage(uuid.NewV4().String()),
	}
}
