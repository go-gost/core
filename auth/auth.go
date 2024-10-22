package auth

import "context"

type Options struct{}

type Option func(opts *Options)

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(ctx context.Context, user, password string, opts ...Option) (id string, ok bool)
}
