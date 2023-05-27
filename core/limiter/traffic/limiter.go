package traffic

import "context"

type ILimiter interface {
	// Wait blocks with the requested n and returns the result value,
	// the returned value is less or equal to n.
	Wait(ctx context.Context, n int) int
	Limit() int
	Set(n int)
}

type ITrafficLimiter interface {
	In(key string) ILimiter
	Out(key string) ILimiter
}
