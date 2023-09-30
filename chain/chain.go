package chain

import (
	"context"
)

type RouteOptions struct {
	Host string
}

type RouteOption func(opts *RouteOptions)

func WithHostRouteOption(host string) RouteOption {
	return func(opts *RouteOptions) {
		opts.Host = host
	}
}

type Chainer interface {
	Route(ctx context.Context, network, address string, opts ...RouteOption) Route
}
