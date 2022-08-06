package eventstroredb

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	esSerializer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"io"
)

type eventStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewEventStore(log logger.Logger, db *esdb.Client) *eventStore {
	return &eventStore{log: log, db: db}
}

func (e *eventStore) SaveEvents(ctx context.Context, streamID string, events []*es.ESEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStore.SaveEvents")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", streamID))

	eventsData := make([]esdb.EventData, 0, len(events))
	for _, event := range events {
		eventsData = append(eventsData, esSerializer.ToEventData(event))
	}

	stream, err := e.db.AppendToStream(ctx, streamID, esdb.AppendToStreamOptions{}, eventsData...)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	e.log.Debugf("SaveEvents stream: %+v", stream)
	return nil
}

func (e *eventStore) LoadEvents(ctx context.Context, streamID string) ([]*es.ESEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStore.LoadEvents")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", streamID))

	stream, err := e.db.ReadStream(ctx, streamID, esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
		From:      esdb.Revision(1),
	}, 100)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer stream.Close()

	events := make([]*es.ESEvent, 0, 100)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		esEvent, err := esSerializer.ToESEventFromRecordedEvent(event.Event)
		if err != nil {
			return nil, err
		}

		events = append(events, esEvent)
	}

	return events, nil
}
