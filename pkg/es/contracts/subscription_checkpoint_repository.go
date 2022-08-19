package contracts

import "context"

type SubscriptionCheckpointRepository interface {
	Load(subscriptionId string, ctx context.Context) (uint64, error)
	Store(subscriptionId string, position uint64, ctx context.Context) error
}
