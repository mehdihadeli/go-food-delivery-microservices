package infrastructure

import (
	"context"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func (ic *infrastructureConfigurator) configKafka(ctx context.Context) (*kafka.Conn, kafkaClient.Producer, error, func()) {

	kafkaConn, err, kafkaConnCleanup := ic.connectKafkaBrokers(ctx)

	if err != nil {
		return nil, nil, errors.Wrap(err, "i.connectKafkaBrokers"), nil
	}

	if ic.cfg.Kafka.InitTopics {
		ic.initKafkaTopics(ctx, kafkaConn)
	}

	kafkaProducer := kafkaClient.NewProducer(ic.log, ic.cfg.Kafka.Brokers)

	return kafkaConn, kafkaProducer, nil, func() {
		_ = kafkaProducer.Close() // nolint:
		kafkaConnCleanup()
	}
}

func (ic *infrastructureConfigurator) connectKafkaBrokers(ctx context.Context) (*kafka.Conn, error, func()) {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, ic.cfg.Kafka)
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewKafkaCon"), nil
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return nil, errors.Wrap(err, "kafkaConn.Brokers"), nil
	}

	ic.log.Infof("kafka connected to brokers: %+v", brokers)

	return kafkaConn, nil, func() {
		_ = kafkaConn.Close() // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) initKafkaTopics(ctx context.Context, kafkaConn *kafka.Conn) {
	controller, err := kafkaConn.Controller()
	if err != nil {
		ic.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	ic.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		ic.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	ic.log.Infof("established new kafka controller connection: %s", controllerURI)

	orderCreateTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.OrderCreate.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.OrderCreate.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.OrderCreate.ReplicationFactor,
	}

	orderCreatedTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.OrderCreated.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.OrderCreated.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.OrderCreated.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		orderCreateTopic,
		orderCreatedTopic,
	); err != nil {
		ic.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	ic.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{orderCreateTopic, orderCreatedTopic})
}
