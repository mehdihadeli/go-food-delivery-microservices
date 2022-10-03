package eventstroredb

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
	appendResult "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/append_result"
	streamName "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_name"
	readPosition "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_position/read_position"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_position/truncatePosition"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_version"
	esErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	attribute2 "go.opentelemetry.io/otel/attribute"
	"math"
)

// https://developers.eventstore.com/clients/grpc/reading-events.html#reading-from-a-stream
// https://developers.eventstore.com/clients/grpc/appending-events.html#append-your-first-event
type eventStoreDbEventStore struct {
	log        logger.Logger
	client     *esdb.Client
	serializer *EsdbSerializer
}

func NewEventStoreDbEventStore(log logger.Logger, client *esdb.Client, serilizer *EsdbSerializer) *eventStoreDbEventStore {
	return &eventStoreDbEventStore{log: log, client: client, serializer: serilizer}
}

func (e *eventStoreDbEventStore) StreamExists(streamName streamName.StreamName, ctx context.Context) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.StreamExists")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	stream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		esdb.ReadStreamOptions{
			Direction: esdb.Backwards,
			From:      esdb.End{}},
		1)
	if err != nil {
		return false, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_StreamExists:ReadStream] error in reading stream"))
	}

	defer stream.Close()

	return stream != nil, nil
}

func (e *eventStoreDbEventStore) AppendEvents(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	events []*models.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.AppendEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	var eventsData []esdb.EventData
	linq.From(events).SelectT(func(s *models.StreamEvent) esdb.EventData {
		data, err := e.serializer.StreamEventToEventData(s)
		if err != nil {
			return *new(esdb.EventData)
		}
		return data
	}).ToSlice(&eventsData)

	var appendEventsResult *appendResult.AppendEventsResult

	res, err := e.client.AppendToStream(
		ctx,
		streamName.String(),
		esdb.AppendToStreamOptions{
			ExpectedRevision: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(expectedVersion),
		},
		eventsData...)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewAppendToStreamError(err, streamName.String()), "[eventStoreDbEventStore_AppendEvents:AppendToStream] error in appending to stream"))
	}

	appendEventsResult = e.serializer.EsdbWriteResultToAppendEventResult(res)

	span.SetAttributes(attribute.Object("AppendEventsResult", appendEventsResult))

	e.log.Infow("[eventStoreDbEventStore_AppendEvents] events append to stream successfully", logger.Fields{"AppendEventsResult": appendEventsResult, "StreamId": streamName.String()})

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) AppendNewEvents(
	streamName streamName.StreamName,
	events []*models.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.AppendNewEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	appendEventsResult, err := e.AppendEvents(streamName, expectedStreamVersion.NoStream, events, ctx)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewAppendToStreamError(err, streamName.String()), "[eventStoreDbEventStore_AppendNewEvents:AppendEvents] error in appending to stream"))
	}

	span.SetAttributes(attribute.Object("AppendNewEvents", appendEventsResult))

	e.log.Infow("[eventStoreDbEventStore_AppendNewEvents] events append to stream successfully", logger.Fields{"AppendEventsResult": appendEventsResult, "StreamId": streamName.String()})

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) ReadEvents(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	readStream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		esdb.ReadStreamOptions{
			Direction:      esdb.Forwards,
			From:           e.serializer.StreamReadPositionToStreamPosition(readPosition),
			ResolveLinkTos: true,
		},
		count)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_ReadEvents:ReadStream] error in reading stream"))
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.EsdbReadStreamToResolvedEvents] error in converting to resolved events"))
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.ResolvedEventsToStreamEvents] error in converting to stream events"))
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEventsWithMaxCount")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEvents(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsFromStart(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEventsFromStart")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEvents(streamName, readPosition.Start, count, ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwards(streamName streamName.StreamName, readPosition readPosition.StreamReadPosition, count uint64, ctx context.Context) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEventsBackwards")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	readStream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		esdb.ReadStreamOptions{
			Direction:      esdb.Backwards,
			From:           e.serializer.StreamReadPositionToStreamPosition(readPosition),
			ResolveLinkTos: true,
		},
		count)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_ReadEventsBackwards:ReadStream] error in reading stream"))
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.EsdbReadStreamToResolvedEvents] error in converting to resolved events"))
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.ResolvedEventsToStreamEvents] error in converting to stream events"))
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEventsBackwards(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsFromEnd(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEventsBackwards(streamName, readPosition.End, count, ctx)
}

func (e *eventStoreDbEventStore) TruncateStream(
	streamName streamName.StreamName,
	truncatePosition truncatePosition.StreamTruncatePosition,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.TruncateStream")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetTruncateBefore(e.serializer.StreamTruncatePositionToInt64(truncatePosition))
	writeResult, err := e.client.SetStreamMetadata(
		ctx,
		streamName.String(),
		esdb.AppendToStreamOptions{
			ExpectedRevision: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(expectedVersion),
		},
		streamMetadata)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewTruncateStreamError(err, streamName.String()), "[eventStoreDbEventStore_TruncateStream:SetStreamMetadata] error in truncating stream"))
	}

	span.SetAttributes(attribute.Object("WriteResult", writeResult))

	e.log.Infow(fmt.Sprintf("[eventStoreDbEventStore.TruncateStream] stream with id %s truncated successfully", streamName.String()), logger.Fields{"WriteResult": writeResult, "StreamId": streamName.String()})

	return e.serializer.EsdbWriteResultToAppendEventResult(writeResult), nil
}

func (e *eventStoreDbEventStore) DeleteStream(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) error {
	ctx, span := tracing.Tracer.Start(ctx, "eventStoreDbEventStore.DeleteStream")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	deleteResult, err := e.client.DeleteStream(
		ctx,
		streamName.String(),
		esdb.DeleteStreamOptions{
			ExpectedRevision: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(expectedVersion),
		})
	if err != nil {
		return tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewDeleteStreamError(err, streamName.String()), "[eventStoreDbEventStore_DeleteStream:DeleteStream] error in deleting stream"))
	}

	span.SetAttributes(attribute.Object("DeleteResult", deleteResult))

	e.log.Infow(fmt.Sprintf("[eventStoreDbEventStore.DeleteStream] stream with id %s deleted successfully", streamName.String()), logger.Fields{"DeleteResult": deleteResult, "StreamId": streamName.String()})

	return nil
}
