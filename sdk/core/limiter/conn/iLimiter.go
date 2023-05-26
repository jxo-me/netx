package conn

type Limiter interface {
	Allow(n int) bool
	Limit() int
}

type IConnLimiter interface {
	Limiter(key string) Limiter
}
