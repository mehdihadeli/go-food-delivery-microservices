package es

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/contracts"
)

type inMemorySubscriptionCheckpointRepository struct {
	checkpoints map[string]uint64
}

func NewInMemorySubscriptionCheckpointRepository() contracts.SubscriptionCheckpointRepository {
	return &inMemorySubscriptionCheckpointRepository{checkpoints: make(map[string]uint64)}
}

func (i inMemorySubscriptionCheckpointRepository) Load(subscriptionId string, ctx context.Context) (uint64, error) {
	checkpoint := i.checkpoints[subscriptionId]
	if checkpoint == 0 {
		return 0, nil
	}
	return checkpoint, nil
}

func (i inMemorySubscriptionCheckpointRepository) Store(
	subscriptionId string,
	position uint64,
	ctx context.Context,
) error {
	i.checkpoints[subscriptionId] = position
	return nil
}
