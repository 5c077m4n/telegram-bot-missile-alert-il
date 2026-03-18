// Package poller to poll requests
package poller

import (
	"context"
	"time"
)

type Fetcher[T any] func() (T, error)

type Poller[T any] struct {
	fetcher  Fetcher[T]
	interval time.Duration
}

func New[T any](fetcher Fetcher[T], interval time.Duration) *Poller[T] {
	return &Poller[T]{fetcher: fetcher, interval: interval}
}

func (p *Poller[T]) Stream(ctx context.Context) (<-chan T, <-chan error) {
	results := make(chan T)
	errors := make(chan error)

	go func() {
		defer close(results)
		defer close(errors)

		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				data, err := p.fetcher()
				if err != nil {
					errors <- err
					continue
				}

				results <- data
			}
		}
	}()

	return results, errors
}
