package eventstroredb

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/append_result"
	streamName "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_name"
	readPosition "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_position/read_position"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_position/truncatePosition"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_version"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math"
)

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
		tracing.TraceErr(span, err)
		return false, err
	}

	defer stream.Close()

	return stream != nil, nil
}

func (e *eventStoreDbEventStore) AppendEvents(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	events []*es.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.AppendEvents")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	var eventsData []esdb.EventData
	linq.From(events).SelectT(func(s *es.StreamEvent) esdb.EventData {
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
		tracing.TraceErr(span, ErrAppendToStream(err, streamName.String()))
		return nil, ErrAppendToStream(err, streamName.String())
	}

	appendEventsResult = e.serializer.EsdbWriteResultToAppendEventResult(res)

	e.log.Infof("AppendEvents result: %s", jsonSerializer.ColoredPrettyPrint(appendEventsResult))

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) AppendNewEvents(
	streamName streamName.StreamName,
	events []*es.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.AppendNewEvents")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	appendEventsResult, err := e.AppendEvents(streamName, expectedStreamVersion.NoStream, events, ctx)
	if err != nil {
		tracing.TraceErr(span, ErrAppendToStream(err, streamName.String()))
		return nil, ErrAppendToStream(err, streamName.String())
	}

	e.log.Infof("AppendNewEvents result: %s", jsonSerializer.ColoredPrettyPrint(appendEventsResult))

	return appendEventsResult, nil
}

func (e *eventStoreDbEventStore) ReadEvents(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	count uint64,
	ctx context.Context,
) ([]*es.StreamEvent, error) {

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
		tracing.TraceErr(span, ErrReadFromStream(err))
		return nil, ErrReadFromStream(err)
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*es.StreamEvent, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsWithMaxCount")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEvents(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsFromStart(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*es.StreamEvent, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsFromStart")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEvents(streamName, readPosition.Start, count, ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwards(streamName streamName.StreamName, readPosition readPosition.StreamReadPosition, count uint64, ctx context.Context) ([]*es.StreamEvent, error) {
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
		tracing.TraceErr(span, ErrReadFromStream(err))
		return nil, ErrReadFromStream(err)
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(readStream)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return events, nil
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*es.StreamEvent, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount")
	span.LogFields(log.String("StreamName", streamName.String()))
	defer span.Finish()

	return e.ReadEventsBackwards(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

func (e *eventStoreDbEventStore) ReadEventsBackwardsFromEnd(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*es.StreamEvent, error) {

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
		tracing.TraceErr(span, ErrTruncateStream(err, streamName.String()))
		return nil, ErrTruncateStream(err, streamName.String())
	}

	e.log.Infof("TruncateStream result: %s", jsonSerializer.ColoredPrettyPrint(writeResult))

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
		tracing.TraceErr(span, ErrTruncateStream(err, streamName.String()))
		return ErrDeleteStream(err, streamName.String())
	}

	e.log.Infof("DeleteStream result: %s", jsonSerializer.ColoredPrettyPrint(deleteResult))

	return nil
}
