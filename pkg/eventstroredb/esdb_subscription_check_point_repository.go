package eventstroredb

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"io"
	"time"
)

type esdbSubscriptionCheckpointRepository struct {
	client        *esdb.Client
	log           logger.Logger
	esdbSerilizer *EsdbSerializer
}

type CheckpointStored struct {
	Position       uint64
	SubscriptionId string
	CheckpointAt   time.Time
}

func NewEsdbSubscriptionCheckpointRepository(client *esdb.Client, logger logger.Logger, esdbSerializer *EsdbSerializer) *esdbSubscriptionCheckpointRepository {
	return &esdbSubscriptionCheckpointRepository{client: client, log: logger, esdbSerilizer: esdbSerializer}
}

func (e *esdbSubscriptionCheckpointRepository) Load(subscriptionId string, ctx context.Context) (uint64, error) {
	streamName := getCheckpointStreamName(subscriptionId)

	stream, err := e.client.ReadStream(
		ctx,
		streamName,
		esdb.ReadStreamOptions{
			Direction: esdb.Backwards,
			From:      esdb.End{},
		}, 1)

	if errors.Is(err, esdb.ErrStreamNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, errors.WrapIf(err, "db.ReadStream")
	}

	event, err := stream.Recv()
	if errors.Is(err, esdb.ErrStreamNotFound) {
		return 0, errors.WrapIf(err, "stream.Recv")
	}
	if errors.Is(err, io.EOF) {
		return 0, nil
	}
	if err != nil {
		return 0, errors.WrapIf(err, "stream.Recv")
	}

	deserialized, _, err := e.esdbSerilizer.Deserialize(event)
	if err != nil {
		return 0, err
	}

	v, ok := deserialized.(*CheckpointStored)
	if !ok {
		return 0, nil
	}

	stream.Close()

	return v.Position, nil
}

func (e *esdbSubscriptionCheckpointRepository) Store(subscriptionId string, position uint64, ctx context.Context) error {
	checkpoint := &CheckpointStored{SubscriptionId: subscriptionId, Position: position, CheckpointAt: time.Now()}
	streamName := getCheckpointStreamName(subscriptionId)
	eventData, err := e.esdbSerilizer.Serialize(checkpoint, nil)
	if err != nil {
		return errors.WrapIf(err, "esdbSerilizer.Serialize")
	}

	_, err = e.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{ExpectedRevision: esdb.StreamExists{}}, *eventData)

	if errors.Is(err, esdb.ErrWrongExpectedStreamRevision) {
		streamMeta := esdb.StreamMetadata{}
		streamMeta.SetMaxCount(1)

		// WrongExpectedVersionException means that stream did not exist
		// Set the checkpoint stream to have at most 1 event
		// using stream metadata $maxCount property
		_, err := e.client.SetStreamMetadata(
			ctx,
			streamName,
			esdb.AppendToStreamOptions{ExpectedRevision: esdb.NoStream{}},
			streamMeta)

		if err != nil {
			return errors.WrapIf(err, "client.SetStreamMetadata")
		}

		// append event again expecting stream to not exist
		_, err = e.client.AppendToStream(ctx, streamName, esdb.AppendToStreamOptions{ExpectedRevision: esdb.NoStream{}}, *eventData)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func getCheckpointStreamName(subscriptionId string) string {
	return fmt.Sprintf("$cehckpoint_stream_%s", subscriptionId)
}
