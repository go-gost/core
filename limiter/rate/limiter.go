package rate

type Limiter interface {
	Allow(n int) bool
	Limit() float64
}

type RateLimiter interface {
	Limiter(key string) Limiter
}
