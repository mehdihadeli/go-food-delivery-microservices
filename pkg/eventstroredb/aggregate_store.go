package eventstroredb

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
	appendResult "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/append_result"
	streamName "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_name"
	readPosition "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_position/read_position"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models/stream_version"
	esErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"
	"reflect"
)

type esdbAggregateStore[T models.IHaveEventSourcedAggregate] struct {
	log        logger.Logger
	eventStore store.EventStore
	serializer *EsdbSerializer
}

func NewEventStoreAggregateStore[T models.IHaveEventSourcedAggregate](log logger.Logger, eventStore store.EventStore, serializer *EsdbSerializer) *esdbAggregateStore[T] {
	return &esdbAggregateStore[T]{log: log, eventStore: eventStore, serializer: serializer}
}

func (a *esdbAggregateStore[T]) StoreWithVersion(aggregate T, metadata metadata.Metadata, expectedVersion expectedStreamVersion.ExpectedStreamVersion, ctx context.Context) (*appendResult.AppendEventsResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "esdbAggregateStore.StoreWithVersion")
	span.SetAttributes(attribute2.String("AggregateID", aggregate.Id().String()))
	defer span.End()

	if len(aggregate.UncommittedEvents()) == 0 {
		a.log.Infow(fmt.Sprintf("[esdbAggregateStore.StoreWithVersion] No events to store for aggregateId %s", aggregate.Id()), logger.Fields{"AggregateID": aggregate.Id()})
		return appendResult.NoOp, nil
	}

	streamId := streamName.For[T](aggregate)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	var streamEvents []*models.StreamEvent

	linq.From(aggregate.UncommittedEvents()).SelectIndexedT(func(i int, domainEvent domain.IDomainEvent) *models.StreamEvent {
		var inInterface map[string]interface{}
		err := jsonSerializer.DecodeWithMapStructure(domainEvent, &inInterface)
		if err != nil {
			return nil
		}
		return a.serializer.DomainEventToStreamEvent(domainEvent, metadata, int64(i)+aggregate.OriginalVersion())
	}).ToSlice(&streamEvents)

	streamAppendResult, err := a.eventStore.AppendEvents(streamId, expectedVersion, streamEvents, ctx)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIff(err, "[esdbAggregateStore_StoreWithVersion:AppendEvents] error in storing aggregate with id {%d}", aggregate.Id()))
	}

	aggregate.MarkUncommittedEventAsCommitted()

	span.SetAttributes(attribute.Object("Aggregate", aggregate))

	a.log.Infow(fmt.Sprintf("[esdbAggregateStore.StoreWithVersion] aggregate with id %d stored successfully", aggregate.Id()), logger.Fields{"Aggregate": aggregate, "StreamId": streamId})

	return streamAppendResult, nil
}

func (a *esdbAggregateStore[T]) Store(aggregate T, metadata metadata.Metadata, ctx context.Context) (*appendResult.AppendEventsResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "esdbAggregateStore.Store")
	defer span.End()

	expectedVersion := expectedStreamVersion.FromInt64(aggregate.OriginalVersion())

	streamAppendResult, err := a.StoreWithVersion(aggregate, metadata, expectedVersion, ctx)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, errors.WrapIff(err, "[esdbAggregateStore_Store:StoreWithVersion] failed to store aggregate with id{%v}", aggregate.Id()))
	}

	return streamAppendResult, nil
}

func (a *esdbAggregateStore[T]) Load(ctx context.Context, aggregateId uuid.UUID) (T, error) {
	ctx, span := tracing.Tracer.Start(ctx, "esdbAggregateStore.Load")
	defer span.End()

	position := readPosition.Start

	return a.LoadWithReadPosition(ctx, aggregateId, position)
}

func (a *esdbAggregateStore[T]) LoadWithReadPosition(ctx context.Context, aggregateId uuid.UUID, position readPosition.StreamReadPosition) (T, error) {
	ctx, span := tracing.Tracer.Start(ctx, "esdbAggregateStore.LoadWithReadPosition")
	span.SetAttributes(attribute2.String("AggregateID", aggregateId.String()))
	defer span.End()

	var typeNameType T
	aggregateInstance := typeMapper.InstancePointerByTypeName(typeMapper.GetFullTypeName(typeNameType))
	aggregate, ok := aggregateInstance.(T)
	if !ok {
		return *new(T), errors.New(fmt.Sprintf("[esdbAggregateStore_LoadWithReadPosition] aggregate is not a %s", typeMapper.GetFullTypeName(typeNameType)))
	}

	method := reflect.ValueOf(aggregate).MethodByName("NewEmptyAggregate")
	if !method.IsValid() {
		return *new(T), tracing.TraceErrFromSpan(span, errors.New("[esdbAggregateStore_LoadWithReadPosition:MethodByName] aggregate does not have a `NewEmptyAggregate` method"))
	}

	method.Call([]reflect.Value{})

	streamId := streamName.ForID[T](aggregateId)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	streamEvents, err := a.getStreamEvents(streamId, position, ctx)
	if errors.Is(err, esdb.ErrStreamNotFound) || len(streamEvents) == 0 {
		return *new(T), tracing.TraceErrFromSpan(span, errors.WithMessage(esErrors.NewAggregateNotFoundError(err, aggregateId), "[esdbAggregateStore.LoadWithReadPosition] error in loading aggregate"))
	}
	if err != nil {
		return *new(T), tracing.TraceErrFromSpan(span, errors.WrapIff(err, "[esdbAggregateStore.LoadWithReadPosition:MethodByName] error in loading aggregate {%s}", aggregateId.String()))
	}

	var meta metadata.Metadata
	var domainEvents []domain.IDomainEvent

	linq.From(streamEvents).Distinct().SelectT(func(streamEvent *models.StreamEvent) domain.IDomainEvent {
		meta = streamEvent.Metadata
		return streamEvent.Event
	}).ToSlice(&domainEvents)

	err = aggregate.LoadFromHistory(domainEvents, meta)
	if err != nil {
		return *new(T), tracing.TraceErrFromSpan(span, err)
	}

	a.log.Infow(fmt.Sprintf("Loaded aggregate with streamId {%s} and aggregateId {%s}",
		streamId.String(),
		aggregateId.String()),
		logger.Fields{"Aggregate": aggregate, "StreamId": streamId.String()})

	span.SetAttributes(attribute.Object("Aggregate", aggregate))

	return aggregate, nil
}

func (a *esdbAggregateStore[T]) Exists(ctx context.Context, aggregateId uuid.UUID) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "esdbAggregateStore.Exists")
	span.SetAttributes(attribute2.String("AggregateID", aggregateId.String()))
	defer span.End()

	streamId := streamName.ForID[T](aggregateId)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	return a.eventStore.StreamExists(streamId, ctx)
}

func (a *esdbAggregateStore[T]) getStreamEvents(streamId streamName.StreamName, position readPosition.StreamReadPosition, ctx context.Context) ([]*models.StreamEvent, error) {
	pageSize := 500
	var streamEvents []*models.StreamEvent

	for true {
		events, err := a.eventStore.ReadEvents(streamId, position, uint64(pageSize), ctx)
		if err != nil {
			return nil, errors.WrapIff(err, "[esdbAggregateStore_getStreamEvents:ReadEvents] failed to read events")
		}
		streamEvents = append(streamEvents, events...)
		if len(events) < pageSize {
			break
		}
		position = readPosition.FromInt64(int64(len(events)) + position.Value())
	}

	return streamEvents, nil
}
