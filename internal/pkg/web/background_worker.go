package web

import (
	"context"
)

type Worker interface {
	Start(ctx context.Context) chan error
	Stop(ctx context.Context) error
}

type (
	ExecutionFunc func(ctx context.Context) error
	StopFunc      func(ctx context.Context) error
)

type BackgroundWorker struct {
	ctx           context.Context
	executionFunc ExecutionFunc
	stopFunc      StopFunc
	cancelFunc    context.CancelFunc
	errChan       chan error
}

func NewBackgroundWorker(executionFunc ExecutionFunc, stopFunc StopFunc) Worker {
	return &BackgroundWorker{executionFunc: executionFunc, stopFunc: stopFunc, errChan: make(chan error)}
}

func (b BackgroundWorker) Start(ctx context.Context) chan error {
	b.ctx, b.cancelFunc = context.WithCancel(ctx)
	go func() {
		if b.executionFunc == nil {
			return
		}

		err := b.executionFunc(b.ctx)
		if err != nil {
			b.cancelFunc()
			b.errChan <- err
		}
	}()
	return b.errChan
}

func (b BackgroundWorker) Stop(ctx context.Context) error {
	if b.executionFunc == nil {
		return nil
	}
	if b.stopFunc != nil {
		return b.stopFunc(ctx)
	}
	if b.cancelFunc != nil {
		b.cancelFunc()
	}

	return nil
}
