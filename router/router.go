// Package router defines the Router interface for OS-level route table
// queries used in TUN (virtual network device) mode.
package router

import (
	"context"
	"net"
)

// Options holds the runtime parameters for a route query.
type Options struct {
	// ID is a unique identifier for the routing table.
	ID string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// IDOption sets the routing table ID.
func IDOption(id string) Option {
	return func(opts *Options) {
		opts.ID = id
	}
}

// Route represents an entry in the system routing table.
type Route struct {
	// Net is the destination network, e.g. "192.168.0.0/16", "172.10.10.0/24".
	// Deprecated: use Dst instead.
	Net *net.IPNet
	// Dst is the destination address or network.
	Dst string
	// Gateway is the gateway address for the destination network.
	Gateway string
}

// Router queries the OS-level route table. This is used in TUN mode where
// GOST acts as a virtual network device and needs to determine how to route
// packets based on the system's routing table.
type Router interface {
	// GetRoute returns the route for the given destination, or nil if no
	// matching route is found.
	GetRoute(ctx context.Context, dst string, opts ...Option) *Route
}
