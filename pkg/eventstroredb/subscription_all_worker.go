package eventstroredb

//
//import (
//	"context"
//	"github.com/EventStore/EventStore-Client-Go/esdb"
//	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
//)
//
//type esdbSubscriptionAllWorker struct {
//	db                 *esdb.Client
//	cfg                *EventStoreConfig
//	log                logger.Logger
//	subscriptionOption *EventStoreDBSubscriptionToAllOptions
//}
//
//type EsdbSubscriptionAllWorker interface {
//	Run() error
//}
//
//type EventStoreDBSubscriptionToAllOptions struct {
//	SubscriptionId              string
//	FilterOptions               *esdb.SubscriptionFilterOptions
//	Credentials                 *esdb.Credentials
//	ResolveLinkTos              bool
//	IgnoreDeserializationErrors bool
//}
//
//func NewEsdbSubscriptionAllWorker(log logger.Logger, db *esdb.Client, cfg *EventStoreConfig, subscriptionOption *EventStoreDBSubscriptionToAllOptions) *esdbSubscriptionAllWorker {
//	if subscriptionOption.SubscriptionId == "" {
//		subscriptionOption.SubscriptionId = "default"
//	}
//
//	if subscriptionOption.FilterOptions == nil {
//		subscriptionOption.FilterOptions.SubscriptionFilter = esdb.ExcludeSystemEventsFilter()
//	}
//
//	return &esdbSubscriptionAllWorker{db: db, cfg: cfg, log: log, subscriptionOption: subscriptionOption}
//}
//
//func (s *esdbSubscriptionAllWorker) Run(ctx context.Context, prefixes []string, poolSize int, worker Worker) error {
//
//	return nil
//}
