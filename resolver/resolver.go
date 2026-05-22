// Package resolver defines the Resolver interface for DNS resolution.
package resolver

import (
	"context"
	"errors"
	"net"
)

var (
	// ErrInvalid is returned when a resolver is not properly configured.
	ErrInvalid = errors.New("invalid resolver")
)

// Options holds the initialization parameters for a Resolver.
type Options struct{}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// Resolver resolves hostnames to IP addresses. It is used by the Router
// to resolve destination addresses before dialing through the proxy chain.
// Implementations may use standard DNS, DNS-over-HTTPS, or custom resolution
// logic.
type Resolver interface {
	// Resolve returns the IPv4 and IPv6 addresses for the given host.
	// The network parameter should be "ip", "ip4", or "ip6"; the default
	// network "ip" returns both IPv4 and IPv6 addresses.
	Resolve(ctx context.Context, network, host string, opts ...Option) ([]net.IP, error)
}
