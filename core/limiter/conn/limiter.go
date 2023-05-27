package conn

type ILimiter interface {
	Allow(n int) bool
	Limit() int
}

type IConnLimiter interface {
	Limiter(key string) ILimiter
}
