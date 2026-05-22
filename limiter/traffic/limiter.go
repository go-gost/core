// Package traffic defines bandwidth/traffic limiter interfaces for
// controlling data transfer rates.
package traffic

import (
	"context"

	"github.com/go-gost/core/limiter"
)

// Limiter controls data transfer bandwidth. It blocks until the requested
// amount of bytes is within the rate limit.
type Limiter interface {
	// Wait blocks with the requested n bytes and returns the actual granted
	// amount, which may be less than or equal to n.
	Wait(ctx context.Context, n int) int
	// Limit returns the current rate limit in bytes per second.
	Limit() int
	// Set updates the rate limit in bytes per second.
	Set(n int)
}

// TrafficLimiter provides per-key traffic limiters for ingress and egress
// traffic. This allows separate bandwidth limits per client, service, or
// other key dimensions.
type TrafficLimiter interface {
	// In returns the Limiter for ingress (download) traffic for the given key.
	In(ctx context.Context, key string, opts ...limiter.Option) Limiter
	// Out returns the Limiter for egress (upload) traffic for the given key.
	Out(ctx context.Context, key string, opts ...limiter.Option) Limiter
}
