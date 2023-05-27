package rate

type ILimiter interface {
	Allow(n int) bool
	Limit() float64
}

type IRateLimiter interface {
	Limiter(key string) ILimiter
}
