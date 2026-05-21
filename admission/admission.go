package admission

import "context"

type Options struct {
	Service string
	Network string
}

type Option func(opts *Options)

func WithService(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

func WithNetwork(network string) Option {
	return func(opts *Options) {
		opts.Network = network
	}
}

type Admission interface {
	Admit(ctx context.Context, network, addr string, opts ...Option) bool
}
