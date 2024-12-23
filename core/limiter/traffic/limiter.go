package traffic

import (
	"context"

	"github.com/jxo-me/netx/core/limiter"
)

type ILimiter interface {
	// Wait blocks with the requested n and returns the result value,
	// the returned value is less or equal to n.
	Wait(ctx context.Context, n int) int
	Limit() int
	Set(n int)
}

type ITrafficLimiter interface {
	In(ctx context.Context, key string, opts ...limiter.Option) ILimiter
	Out(ctx context.Context, key string, opts ...limiter.Option) ILimiter
}
