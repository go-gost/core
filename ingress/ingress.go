package ingress

import "context"

type Options struct{}

type Option func(opts *Options)

type Rule struct {
	// Hostname is the hostname match pattern, e.g. example.com, *.example.org or .example.com.
	Hostname string
	// Endpoint is the tunnel ID for the hostname.
	Endpoint string
}

type Ingress interface {
	// SetRule adds or updates a rule for the ingress.
	SetRule(ctx context.Context, rule *Rule, opts ...Option) bool
	// GetRule queries a rule by host.
	GetRule(ctx context.Context, host string, opts ...Option) *Rule
}
