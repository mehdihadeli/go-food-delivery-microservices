package web

import "context"

type WorkersRunner struct {
	workers []Worker
	errChan chan error
}

func NewWorkersRunner(workers []Worker) *WorkersRunner {
	return &WorkersRunner{workers: workers, errChan: make(chan error)}
}

func (r *WorkersRunner) Start(ctx context.Context) chan error {
	if r.workers == nil || len(r.workers) == 0 {
		return nil
	}

	for _, w := range r.workers {
		err := w.Start(ctx)
		go func() {
			for {
				select {
				case e := <-err:
					r.errChan <- e
					return
				case <-ctx.Done():
					stopErr := r.Stop(ctx)
					if stopErr != nil {
						r.errChan <- stopErr
						return
					}
					return
				}
			}
		}()
	}

	return r.errChan
}

func (r *WorkersRunner) Stop(ctx context.Context) error {
	if r.workers == nil || len(r.workers) == 0 {
		return nil
	}

	for _, w := range r.workers {
		err := w.Stop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
