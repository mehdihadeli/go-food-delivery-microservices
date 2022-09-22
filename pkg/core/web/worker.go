package web

import "context"

type Worker interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
