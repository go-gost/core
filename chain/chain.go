// Package chain defines the proxy chain interfaces — the core routing layer of GOST.
// A proxy chain consists of multiple nodes through which traffic is forwarded
// sequentially to reach the destination.
package chain

import (
	"context"
)

// RouteOptions holds the runtime parameters for route selection.
type RouteOptions struct {
	// Host is the target hostname, used for host-based routing decisions.
	Host string
}

// RouteOption is a functional option for configuring RouteOptions.
type RouteOption func(opts *RouteOptions)

// WithHostRouteOption sets the target hostname.
func WithHostRouteOption(host string) RouteOption {
	return func(opts *RouteOptions) {
		opts.Host = host
	}
}

// Chainer selects a Route for the given target address. It represents a
// collection of proxy chains, from which the best matching route is chosen.
type Chainer interface {
	// Route returns a Route for the target network and address, or nil if
	// no suitable route is available.
	Route(ctx context.Context, network, address string, opts ...RouteOption) Route
}
