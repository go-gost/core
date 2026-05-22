// Package bypass defines the Bypass interface for deciding whether a target
// address should skip the proxy chain and connect directly.
package bypass

import (
	"context"
)

// Options holds the initialization parameters for a Bypass rule.
type Options struct {
	// Service is the service name this bypass rule belongs to.
	Service string
	// Host is the HTTP Host header value.
	Host string
	// Path is the HTTP request path.
	Path string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// WithService sets the service name.
func WithService(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// WithHostOption sets the HTTP Host header.
func WithHostOption(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

// WithPathOption sets the HTTP request path.
func WithPathOption(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

// Bypass is an address filter that determines whether a target address should
// bypass the proxy chain and be reached directly.
// When Contains returns true for an address, that address will not go through
// the proxy forwarding chain.
type Bypass interface {
	// IsWhitelist reports whether this rule operates as a whitelist (true)
	// or blacklist (false). In whitelist mode, only matched addresses bypass
	// the proxy; in blacklist mode, matched addresses go through the proxy.
	IsWhitelist() bool
	// Contains reports whether the given address matches this bypass rule.
	Contains(ctx context.Context, network, addr string, opts ...Option) bool
}
