package kafka

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
	"wb-tech-l0/internal/domain/app/ports"
)

type Reader interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Worker struct {
	log      ports.Logger
	handlers []Reader
}

func NewWorker(
	log ports.Logger,
	handlers ...Reader,
) *Worker {
	return &Worker{
		log:      log,
		handlers: handlers,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("starting kafka reader worker", "total_readers", len(w.handlers))

	errChan := make(chan error, 1)
	errGroup, ctx := errgroup.WithContext(ctx)
	go func() {
		errChan <- errGroup.Wait()
	}()

	for idx, handler := range w.handlers {
		func(idx int, handler Reader) {
			errGroup.Go(func() error {
				return handler.Start(ctx)
			})
		}(idx, handler)
	}

	select {
	case err := <-errChan:
		w.log.Error("kafka reader received critical error, initiating shutdown", "error", err.Error())
		return fmt.Errorf("reader error: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (w *Worker) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	go func() {
		defer close(done)
		for i := len(w.handlers) - 1; i >= 0; i-- {
			if err := w.handlers[i].Stop(ctx); err != nil {
				w.log.Error("kafka reader stopped with error", "error", err.Error())
			}
		}
	}()

	select {
	case <-done:
		w.log.Info("all readers stopped")
		return nil
	case <-ctx.Done():
		w.log.Warn("forced shutdown due to context timeout")
		return ctx.Err()
	}
}
