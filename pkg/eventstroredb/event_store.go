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
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math"
)

//https://developers.eventstore.com/clients/grpc/reading-events.html#reading-from-a-stream
//https://developers.eventstore.com/clients/grpc/appending-events.html#append-your-first-event
type eventStoreDbEventStore struct {
	log        logger.Logger
	client     *esdb.Client
	serializer *EsdbSerializer
}

func NewEventStoreDbEventStore(log logger.Logger, client *esdb.Client, serilizer *EsdbSerializer) *eventStoreDbEventStore {
	return &eventStoreDbEventStore{log: log, client: client, serializer: serilizer}
}

func (e *eventStoreDbEventStore) StreamExists(streamName streamName.StreamName, ctx context.Context) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.StreamExists")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	stream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		esdb.ReadStreamOptions{
			Direction: esdb.Backwards,
			From:      esdb.End{}},
		1)
	if err != nil {
		return false, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_StreamExists:ReadStream] error in reading stream"))
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.AppendEvents")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

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
		return nil, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewAppendToStreamError(err, streamName.String()), "[eventStoreDbEventStore_AppendEvents:AppendToStream] error in appending to stream"))
	}

	appendEventsResult = e.serializer.EsdbWriteResultToAppendEventResult(res)

	span.LogFields(log.Object("AppendEventsResult", appendEventsResult))

	e.log.Infow("[eventStoreDbEventStore_AppendEvents] events append to stream successfully", logger.Fields{"AppendEventsResult": appendEventsResult, "StreamId": streamName.String()})

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) AppendNewEvents(
	streamName streamName.StreamName,
	events []*models.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.AppendNewEvents")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	appendEventsResult, err := e.AppendEvents(streamName, expectedStreamVersion.NoStream, events, ctx)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewAppendToStreamError(err, streamName.String()), "[eventStoreDbEventStore_AppendNewEvents:AppendEvents] error in appending to stream"))
	}

	span.LogFields(log.Object("AppendNewEvents", appendEventsResult))

	e.log.Infow("[eventStoreDbEventStore_AppendNewEvents] events append to stream successfully", logger.Fields{"AppendEventsResult": appendEventsResult, "StreamId": streamName.String()})

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) ReadEvents(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEvents")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

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
		return nil, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_ReadEvents:ReadStream] error in reading stream"))
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.EsdbReadStreamToResolvedEvents] error in converting to resolved events"))
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.ResolvedEventsToStreamEvents] error in converting to stream events"))
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsWithMaxCount")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEvents(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsFromStart(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsFromStart")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEvents(streamName, readPosition.Start, count, ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwards(streamName streamName.StreamName, readPosition readPosition.StreamReadPosition, count uint64, ctx context.Context) ([]*models.StreamEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsBackwards")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

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
		return nil, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewReadStreamError(err), "[eventStoreDbEventStore_ReadEventsBackwards:ReadStream] error in reading stream"))
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.EsdbReadStreamToResolvedEvents] error in converting to resolved events"))
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, tracing.TraceWithErr(span, errors.WrapIf(err, "[eventStoreDbEventStore_ReadEvents.ResolvedEventsToStreamEvents] error in converting to stream events"))
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEventsBackwards(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsFromEnd(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsBackwardsFromEnd")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEventsBackwards(streamName, readPosition.End, count, ctx)
}

func (e *eventStoreDbEventStore) TruncateStream(
	streamName streamName.StreamName,
	truncatePosition truncatePosition.StreamTruncatePosition,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.TruncateStream")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

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
		return nil, tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewTruncateStreamError(err, streamName.String()), "[eventStoreDbEventStore_TruncateStream:SetStreamMetadata] error in truncating stream"))
	}

	span.LogFields(log.Object("WriteResult", writeResult))

	e.log.Infow(fmt.Sprintf("[eventStoreDbEventStore.TruncateStream] stream with id %s truncated successfully", streamName.String()), logger.Fields{"WriteResult": writeResult, "StreamId": streamName.String()})

	return e.serializer.EsdbWriteResultToAppendEventResult(writeResult), nil
}

func (e *eventStoreDbEventStore) DeleteStream(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.DeleteStream")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	deleteResult, err := e.client.DeleteStream(
		ctx,
		streamName.String(),
		esdb.DeleteStreamOptions{
			ExpectedRevision: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(expectedVersion),
		})
	if err != nil {
		return tracing.TraceWithErr(span, errors.WithMessage(esErrors.NewDeleteStreamError(err, streamName.String()), "[eventStoreDbEventStore_DeleteStream:DeleteStream] error in deleting stream"))
	}

	span.LogFields(log.Object("DeleteResult", deleteResult))

	e.log.Infow(fmt.Sprintf("[eventStoreDbEventStore.DeleteStream] stream with id %s deleted successfully", streamName.String()), logger.Fields{"DeleteResult": deleteResult, "StreamId": streamName.String()})

	return nil
}
