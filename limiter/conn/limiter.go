// Package conn defines connection-level limiter interfaces for controlling
// the number of active connections.
package conn

// Limiter controls the number of concurrent operations. It reports whether
// a given count of new connections is allowed without exceeding the limit.
type Limiter interface {
	// Allow reports whether n connections can be accepted.
	Allow(n int) bool
	// Limit returns the maximum allowed count.
	Limit() int
}

// ConnLimiter provides per-key connection limiters. Keys are typically
// client IPs or service names, allowing different limits for different
// sources or services.
type ConnLimiter interface {
	// Limiter returns the Limiter for the given key.
	Limiter(key string) Limiter
}
