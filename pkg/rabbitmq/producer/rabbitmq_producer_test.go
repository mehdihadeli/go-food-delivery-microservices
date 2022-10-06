package producer

//func Test_Publish_Message(t *testing.T) {
//	test.SkipCI(t)
//	ctx := context.Background()
//	tp, err := tracing.AddOtelTracing(&otel.OpenTelemetryConfig{ServiceName: "test", Enabled: true, AlwaysOnSampler: true, JaegerExporterConfig: &otel.JaegerExporterConfig{AgentHost: "localhost", AgentPort: "6831"}})
//	if err != nil {
//		return
//	}
//	defer tp.Shutdown(ctx)
//
//	conn, err := types.NewRabbitMQConnection(context.Background(), &config.RabbitMQConfig{
//		RabbitMqHostOptions: &config.RabbitMqHostOptions{
//			UserName: "guest",
//			Password: "guest",
//			HostName: "localhost",
//			Port:     5672,
//		},
//	})
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//
//	rabbitmqProducer, err := NewRabbitMQProducer(conn,
//		func(builder *configurations.rabbitMQProducerConfigurationBuilder) {
//			builder.WithExchangeType(types.ExchangeTopic)
//		}, defaultLogger.Logger, json.NewJsonEventSerializer())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = rabbitmqProducer.PublishMessage(context.Background(), NewProducerMessage("test"), nil)
//	if err != nil {
//		return
//	}
//}
//
//type ProducerMessage struct {
//	*types2.Message
//	Data string
//}
//
//func NewProducerMessage(data string) *ProducerMessage {
//	return &ProducerMessage{
//		Data:    data,
//		Message: types2.NewMessage(uuid.NewV4().String()),
//	}
//}
