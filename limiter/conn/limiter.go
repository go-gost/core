package conn

type Limiter interface {
	Allow(n int) bool
	Limit() int
}

type ConnLimiter interface {
	Limiter(key string) Limiter
}
