// Package ingress defines the Ingress interface for hostname-based tunnel
// routing. It maps incoming hostnames to tunnel endpoints, enabling reverse
// proxy scenarios where external traffic is routed to internal services.
package ingress

import "context"

// Options holds the initialization parameters for an Ingress.
type Options struct {
	// Service is the service name this ingress rule belongs to.
	Service string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// WithService sets the service name.
func WithService(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// Rule defines a mapping from a hostname pattern to a tunnel endpoint.
type Rule struct {
	// Hostname is the hostname match pattern, e.g. "example.com",
	// "*.example.org", or ".example.com".
	Hostname string
	// Endpoint is the tunnel ID that traffic for the matching hostname
	// should be routed to.
	Endpoint string
}

// Ingress manages hostname-to-endpoint routing rules for incoming traffic.
// It maps domain names to tunnel endpoints, enabling multi-tenant reverse
// proxy setups where different hostnames are serviced by different backends.
type Ingress interface {
	// SetRule adds or updates a routing rule.
	SetRule(ctx context.Context, rule *Rule, opts ...Option) bool
	// GetRule queries a routing rule by hostname.
	GetRule(ctx context.Context, host string, opts ...Option) *Rule
}
