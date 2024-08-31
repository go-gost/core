package traffic

import (
	"context"

	"github.com/go-gost/core/limiter"
)

type Limiter interface {
	// Wait blocks with the requested n and returns the result value,
	// the returned value is less or equal to n.
	Wait(ctx context.Context, n int) int
	Limit() int
	Set(n int)
}

type TrafficLimiter interface {
	In(ctx context.Context, key string, opts ...limiter.Option) Limiter
	Out(ctx context.Context, key string, opts ...limiter.Option) Limiter
}
