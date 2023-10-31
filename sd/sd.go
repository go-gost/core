package sd

import (
	"context"
)

type Options struct{}

type Option func(opts *Options)

type Service struct {
	Node    string
	Name    string
	Network string
	Address string
}

// SD is the service discovery interface.
type SD interface {
	Register(ctx context.Context, name string, network string, address string, opts ...Option) error
	Deregister(ctx context.Context, name string) error
	Renew(ctx context.Context, name string) error
	Get(ctx context.Context, name string) ([]*Service, error)
}
