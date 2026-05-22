// Package auth defines the Authenticator interface for verifying client identities.
package auth

import "context"

// Options holds the initialization parameters for an Authenticator.
type Options struct {
	// Service is the service name this authenticator belongs to.
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

// Authenticator is the user authentication interface.
// In proxy scenarios, authentication happens after the client connection is
// established and before any request is processed. Only authenticated connections
// are allowed to proceed.
type Authenticator interface {
	// Authenticate validates the username and password. It returns the user
	// identifier id and whether authentication succeeded.
	// The id can be used for subsequent authorization and audit trails.
	Authenticate(ctx context.Context, user, password string, opts ...Option) (id string, ok bool)
}
