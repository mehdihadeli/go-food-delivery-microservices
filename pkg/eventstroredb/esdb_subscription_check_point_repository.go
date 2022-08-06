package eventstroredb

//
//type esdbSubscriptionCheckpointRepository struct {
//	client *esdb.Client
//	log    logger.Logger
//}
//
//func NewesdbSubscriptionCheckpointRepository(client *esdb.Client, logger logger.Logger) *esdbSubscriptionCheckpointRepository {
//	return &esdbSubscriptionCheckpointRepository{client: client, log: logger}
//}
//
//func (e *esdbSubscriptionCheckpointRepository) Load(subscriptionId string, ctx context.Context) (int, error) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "esdbSubscriptionCheckpointRepository.Load")
//	streamName := getCheckpointStreamName(subscriptionId)
//
//	stream, err := e.client.ReadStream(
//		ctx,
//		streamName,
//		esdb.ReadStreamOptions{
//			Direction: esdb.Backwards,
//			From:      esdb.End{},
//		}, 1)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		return 0, errors.Wrap(err, "db.ReadStream")
//	}
//	defer stream.Close()
//
//	for {
//		event, err := stream.Recv()
//		if errors.Is(err, esdb.ErrStreamNotFound) {
//			tracing.TraceErr(span, err)
//			return 0, errors.Wrap(err, "stream.Recv")
//		}
//		if errors.Is(err, io.EOF) {
//			break
//		}
//		if err != nil {
//			tracing.TraceErr(span, err)
//			return 0, errors.Wrap(err, "stream.Recv")
//		}
//
//		esEvent, err := esSerializer.ToESEventFromRecordedEvent(event.Event)
//		if err != nil {
//			tracing.TraceErr(span, err)
//			return 0, errors.Wrap(err, "serializer.ToESEventFromRecordedEvent")
//		}
//		if err := aggregate.RaiseEvent(esEvent); err != nil {
//			tracing.TraceErr(span, err)
//			return 0, errors.Wrap(err, "RaiseEvent")
//		}
//		e.log.Debugf("(Load) esEvent: {%s}", esEvent.String())
//	}
//}
//
//func (e *esdbSubscriptionCheckpointRepository) Store(subscriptionId string, position int, ctx context.Context) {
//
//}
//
//func getCheckpointStreamName(subscriptionId string) string {
//	return fmt.Sprintf("$cehckpoint_%s", subscriptionId)
//}
