package selector

import (
	"context"
	"time"

	"github.com/jxo-me/netx/core/selector"
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
	strategy selector.IStrategy[T]
	filters  []selector.IFilter[T]
}

func NewSelector[T any](strategy selector.IStrategy[T], filters ...selector.IFilter[T]) selector.ISelector[T] {
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
