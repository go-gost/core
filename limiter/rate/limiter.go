// Package rate defines rate-limiting interfaces for controlling the rate
// of new connections or requests.
package rate

// Limiter controls the rate of operations. It reports whether a given number
// of operations is within the allowed rate.
type Limiter interface {
	// Allow reports whether n operations can proceed at the current rate.
	Allow(n int) bool
	// Limit returns the current rate limit (operations per second).
	Limit() float64
}

// RateLimiter provides per-key rate limiters. Keys are typically client IPs
// or service names, allowing different rate limits for different sources.
type RateLimiter interface {
	// Limiter returns the Limiter for the given key.
	Limiter(key string) Limiter
}
