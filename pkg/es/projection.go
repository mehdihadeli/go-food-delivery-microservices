package es

import (
	"context"
)

// Projection When method works and process Event's like Aggregate's for interacting with read database.
type Projection interface {
	When(ctx context.Context, evt Event) error
}
