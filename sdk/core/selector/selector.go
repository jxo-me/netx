package selector

import (
	"context"
	"time"
)

// default options for FailFilter
const (
	DefaultMaxFails    = 1
	DefaultFailTimeout = 10 * time.Second
)

const (
	labelWeight      = "weight"
	labelBackup      = "backup"
	labelMaxFails    = "maxFails"
	labelFailTimeout = "failTimeout"
)

type defaultSelector[T any] struct {
	strategy Strategy[T]
	filters  []Filter[T]
}

func NewSelector[T any](strategy Strategy[T], filters ...Filter[T]) Selector[T] {
	return &defaultSelector[T]{
		filters:  filters,
		strategy: strategy,
	}
}

func (s *defaultSelector[T]) Select(ctx context.Context, vs ...T) (v T) {
	for _, filter := range s.filters {
		vs = filter.Filter(ctx, vs...)
	}
	if len(vs) == 0 {
		return
	}
	return s.strategy.Apply(ctx, vs...)
}
