package contracts

import "context"

type SubscriptionCheckpointRepository interface {
	Load(subscriptionId string, ctx context.Context) int
	Store(subscriptionId string, position int, ctx context.Context)
}
