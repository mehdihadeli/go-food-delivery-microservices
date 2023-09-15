package health

import "context"

type Health interface {
	CheckHealth(ctx context.Context) error
	GetHealthName() string
}
