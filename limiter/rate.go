package limiter

type Limiter interface {
	// Limit checks the requested size b and returns the limit size,
	// the returned value is less or equal to b.
	Limit(b int) int
}

type RateLimiter interface {
	Input() Limiter
	Output() Limiter
}
