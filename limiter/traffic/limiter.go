package traffic

import "context"

type Limiter interface {
	// Wait blocks with the requested n and returns the result value,
	// the returned value is less or equal to n.
	Wait(ctx context.Context, n int) int
	Limit() int
	Set(n int)
}

type TrafficLimiter interface {
	In(key string) Limiter
	Out(key string) Limiter
}
