package eventstroredb

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/goccy/go-reflect"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	esSerializer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_name"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"math"
)

const (
	count = math.MaxInt64
)

type aggregateStore[T contracts.IEventSourcedAggregateRoot] struct {
	log logger.Logger
	db  *esdb.Client
}

func NewEventStoreAggregateStore[T contracts.IEventSourcedAggregateRoot](log logger.Logger, db *esdb.Client) *aggregateStore[T] {
	return &aggregateStore[T]{log: log, db: db}
}

func (a *aggregateStore[T]) Load(ctx context.Context, aggregateId uuid.UUID) (T, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Load")
	defer span.Finish()

	var typeNameType T
	aggregateInstance := typeMapper.TypePointerInstanceByName(typeMapper.GetTypeName(typeNameType))
	aggregate, ok := aggregateInstance.(T)
	if !ok {
		return *new(T), errors.New("aggregateStore.Load: aggregate is not a T")
	}

	method := reflect.ValueOf(aggregate).MethodByName("NewEmptyAggregate")
	if !method.IsValid() {
		return *new(T), errors.New("aggregateStore.Load: aggregate does not have a NewEmptyAggregate method")
	}

	method.Call([]reflect.Value{})

	streamId := streamName.ForID[T](aggregateId)
	span.LogFields(log.String("AggregateId", aggregateId.String()))
	span.LogFields(log.String("StreamId", streamId))

	stream, err := a.db.ReadStream(ctx, streamId, esdb.ReadStreamOptions{}, count)
	if err != nil {
		tracing.TraceErr(span, err)
		return *new(T), errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	for {
		event, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			tracing.TraceErr(span, err)
			return *new(T), errors.Wrap(err, "stream.Recv")
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return *new(T), errors.Wrap(err, "stream.Recv")
		}

		esEvent, err := esSerializer.ToESEventFromRecordedEvent(event.Event)

		deserializedEvent, err := esSerializer.DeserializeEventFromESEvent(esEvent)
		if err != nil {
			a.log.Errorf("(loadEvents) serializer.DeserializeEvent err: %v", err)
			return *new(T), tracing.TraceWithErr(span, errors.Wrap(err, "serializer.DeserializeEvent"))
		}
		a.log.Debugf("(loadEvents) deserializedEvents: %v", serializer.ColoredPrettyPrint(deserializedEvent))

		deserializedMeta, err := esSerializer.DeserializeMetadataFromESEvent(esEvent)
		if err != nil {
			return *new(T), err
		}
		a.log.Debugf("(loadMeta) deserializedMeta: %v", serializer.ColoredPrettyPrint(deserializedMeta))

		if err := aggregate.RaiseEvent(deserializedEvent); err != nil {
			tracing.TraceErr(span, err)
			return *new(T), errors.Wrap(err, "RaiseEvent")
		}
		a.log.Debugf("(Load) esEvent: {%s}", esEvent.String())
	}

	a.log.Debugf("(Load) aggregate: {%s}", aggregate.String())

	return aggregate, nil
}

func (a *aggregateStore[T]) Store(ctx context.Context, aggregate T, metadata *domain.Metadata) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Save")
	defer span.Finish()
	span.LogFields(log.String("aggregate", aggregate.String()))

	if len(aggregate.GetUncommittedEvents()) == 0 {
		a.log.Debugf("(Save) [no uncommittedEvents] len: {%d}", len(aggregate.GetUncommittedEvents()))
		return nil
	}

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {

		event, err := esSerializer.SerializeToESEvent(aggregate, event, metadata)
		if err != nil {
			a.log.Errorf("(Save) serializer.SerializeEvent err: %v", err)
			return tracing.TraceWithErr(span, errors.Wrap(err, "serializer.SerializeEvent"))
		}

		eventsData = append(eventsData, esSerializer.ToEventData(event))
	}

	// check for aggregate.GetVersion() == 0 or len(aggregate.GetAppliedEvents()) == 0 means new aggregate
	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() == 0 {
		expectedRevision = esdb.NoStream{}
		a.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

		appendStream, err := a.db.AppendToStream(
			ctx,
			streamName.For(aggregate),
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "db.AppendToStream")
		}

		a.log.Debugf("(Save) stream: {%+v}", appendStream)
		return nil
	}

	readOps := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := a.db.ReadStream(context.Background(), streamName.For(aggregate), readOps, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "stream.Recv")
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	a.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

	appendStream, err := a.db.AppendToStream(
		ctx,
		streamName.For(aggregate),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "db.AppendToStream")
	}

	a.log.Debugf("(Save) stream: {%+v}", appendStream)
	aggregate.MarkUncommittedEventAsCommitted()
	return nil
}

func (a *aggregateStore[T]) Exists(ctx context.Context, aggregateId uuid.UUID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Exists")
	defer span.Finish()
	streamId := streamName.ForID[T](aggregateId)

	span.LogFields(log.String("AggregateId", aggregateId.String()))
	span.LogFields(log.String("StreamId", streamId))

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := a.db.ReadStream(ctx, streamId, readStreamOptions, 1)
	if err != nil {
		return false, errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	for {
		_, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			tracing.TraceErr(span, err)
			return false, errors.Wrap(esdb.ErrStreamNotFound, "stream.Recv")
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return false, errors.Wrap(err, "stream.Recv")
		}
	}

	return true, nil
}
