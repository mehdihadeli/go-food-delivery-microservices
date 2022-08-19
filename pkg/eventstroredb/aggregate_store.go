package eventstroredb

import (
	"context"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	appendResult "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/append_result"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_name"
	readPosition "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_position/read_position"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_version"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"reflect"
)

type esdbAggregateStore[T es.IHaveEventSourcedAggregate] struct {
	log        logger.Logger
	eventStore store.EventStore
	serializer *EsdbSerializer
}

func NewEventStoreAggregateStore[T es.IHaveEventSourcedAggregate](log logger.Logger, eventStore store.EventStore, serializer *EsdbSerializer) *esdbAggregateStore[T] {
	return &esdbAggregateStore[T]{log: log, eventStore: eventStore, serializer: serializer}
}

func (a *esdbAggregateStore[T]) StoreWithVersion(aggregate T, metadata *core.Metadata, expectedVersion expectedStreamVersion.ExpectedStreamVersion, ctx context.Context) (*appendResult.AppendEventsResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbAggregateStore.StoreWithVersion")
	defer span.Finish()
	span.LogFields(log.String("Aggregate", jsonSerializer.ColoredPrettyPrint(aggregate)))

	if len(aggregate.UncommittedEvents()) == 0 {
		a.log.Debugf("No events to store for aggregateId %s", aggregate.Id())
		return appendResult.NoOp, nil
	}

	streamId := streamName.For[T](aggregate)
	span.LogFields(log.String("StreamId", streamId.String()))

	var streamEvents []*es.StreamEvent

	linq.From(aggregate.UncommittedEvents()).SelectIndexedT(func(i int, domainEvent domain.IDomainEvent) *es.StreamEvent {
		var inInterface map[string]interface{}
		err := jsonSerializer.DecodeWithMapStructure(domainEvent, &inInterface)
		if err != nil {
			return nil
		}
		return a.serializer.DomainEventToStreamEvent(domainEvent, metadata, int64(i)+aggregate.OriginalVersion())
	}).ToSlice(&streamEvents)

	streamAppendResult, err := a.eventStore.AppendEvents(streamId, expectedVersion, streamEvents, ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrapf(err, "failed to store aggregate {%s}", jsonSerializer.ColoredPrettyPrint(aggregate))
	}

	a.log.Debugf("StreamAppendResult for aggregateId %s: %s", aggregate.Id(), jsonSerializer.ColoredPrettyPrint(streamAppendResult))
	aggregate.MarkUncommittedEventAsCommitted()

	return streamAppendResult, nil
}

func (a *esdbAggregateStore[T]) Store(aggregate T, metadata *core.Metadata, ctx context.Context) (*appendResult.AppendEventsResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbAggregateStore.Store")
	defer span.Finish()
	span.LogFields(log.String("Aggregate", jsonSerializer.ColoredPrettyPrint(aggregate)))

	expectedVersion := expectedStreamVersion.FromInt64(aggregate.OriginalVersion())

	streamAppendResult, err := a.StoreWithVersion(aggregate, metadata, expectedVersion, ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrapf(err, "failed to store aggregate {%s}", jsonSerializer.ColoredPrettyPrint(aggregate))
	}

	a.log.Debugf("StreamAppendResult for aggregateId %s: %s", aggregate.Id(), jsonSerializer.ColoredPrettyPrint(streamAppendResult))

	return streamAppendResult, nil
}

func (a *esdbAggregateStore[T]) Load(ctx context.Context, aggregateId uuid.UUID) (T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbAggregateStore.Load")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", aggregateId.String()))

	position := readPosition.Start

	return a.LoadWithReadPosition(ctx, aggregateId, position)
}

func (a *esdbAggregateStore[T]) LoadWithReadPosition(ctx context.Context, aggregateId uuid.UUID, position readPosition.StreamReadPosition) (T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbAggregateStore.LoadWithReadPosition")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", aggregateId.String()))

	var typeNameType T
	aggregateInstance := typeMapper.InstancePointerByTypeName(typeMapper.GetTypeName(typeNameType))
	aggregate, ok := aggregateInstance.(T)
	if !ok {
		return *new(T), errors.New(fmt.Sprintf("aggregate is not a %s", typeMapper.GetTypeName(typeNameType)))
	}

	method := reflect.ValueOf(aggregate).MethodByName("NewEmptyAggregate")
	if !method.IsValid() {
		return *new(T), errors.New("aggregate does not have a `NewEmptyAggregate` method")
	}

	method.Call([]reflect.Value{})

	streamId := streamName.ForID[T](aggregateId)
	span.LogFields(log.String("StreamId", streamId.String()))

	streamEvents, err := a.getStreamEvents(streamId, position, ctx)
	if errors.Is(err, ErrStreamNotFound(err)) || len(streamEvents) == 0 {
		tracing.TraceErr(span, ErrAggregateNotFound(err, aggregateId.String()))
		return *new(T), ErrAggregateNotFound(err, aggregateId.String())
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return *new(T), errors.Wrapf(err, "failed to load aggregate {%s}", aggregateId.String())
	}

	var metadata *core.Metadata
	var domainEvents []domain.IDomainEvent

	linq.From(streamEvents).Distinct().SelectT(func(streamEvent *es.StreamEvent) domain.IDomainEvent {
		metadata = streamEvent.Metadata
		return streamEvent.Event
	}).ToSlice(&domainEvents)

	err = aggregate.LoadFromHistory(domainEvents, metadata)
	if err != nil {
		tracing.TraceErr(span, err)
		return *new(T), err
	}

	a.log.Debugf("Loaded aggregate {%s} from streamId {%s}", jsonSerializer.ColoredPrettyPrint(aggregate), streamId.String())

	return aggregate, nil
}

func (a *esdbAggregateStore[T]) Exists(ctx context.Context, aggregateId uuid.UUID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbAggregateStore.Exists")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", aggregateId.String()))

	streamId := streamName.ForID[T](aggregateId)
	span.LogFields(log.String("StreamId", streamId.String()))

	return a.eventStore.StreamExists(streamId, ctx)
}

func (a *esdbAggregateStore[T]) getStreamEvents(streamId streamName.StreamName, position readPosition.StreamReadPosition, ctx context.Context) ([]*es.StreamEvent, error) {
	pageSize := 500
	var streamEvents []*es.StreamEvent

	for true {
		events, err := a.eventStore.ReadEvents(streamId, position, uint64(pageSize), ctx)
		if err != nil {
			return nil, err
		}
		streamEvents = append(streamEvents, events...)
		if len(events) < pageSize {
			break
		}
		position = readPosition.FromInt64(int64(len(events)) + position.Value())
	}

	return streamEvents, nil
}
