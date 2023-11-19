package router

import (
	"context"
	"net"
)

type Options struct{}

type Option func(opts *Options)

type Route struct {
	// Net is the destination network, e.g. 192.168.0.0/16, 172.10.10.0/24.
	Net *net.IPNet
	// Gateway is the gateway for the destination network.
	Gateway net.IP
}

type Router interface {
	// GetRoute queries a route by destination IP address.
	GetRoute(ctx context.Context, dst net.IP, opts ...Option) *Route
}
