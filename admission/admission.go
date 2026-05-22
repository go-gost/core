// Package admission defines the Admission interface for connection admission control,
// determining whether to allow or deny incoming connections.
package admission

import "context"

// Options holds the initialization parameters for an Admission.
type Options struct {
	// Service is the service name this admission rule belongs to.
	Service string
	// Network is the connection network type, e.g. "tcp", "udp".
	Network string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// WithService sets the service name.
func WithService(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// WithNetwork sets the network type.
func WithNetwork(network string) Option {
	return func(opts *Options) {
		opts.Network = network
	}
}

// Admission is the admission control interface. It decides whether a connection
// to the given address should be allowed or denied.
// Common implementations include IP whitelist/blacklist and domain-based rules.
type Admission interface {
	// Admit reports whether the connection to the given network address is allowed.
	// It returns true if the connection is permitted.
	Admit(ctx context.Context, network, addr string, opts ...Option) bool
}
