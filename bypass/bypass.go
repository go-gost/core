package bypass

import (
	"context"
)

type Options struct {
	Host string
	Path string
}

type Option func(opts *Options)

func WithHostOpton(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

func WithPathOption(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

// Bypass is a filter of address (IP or domain).
type Bypass interface {
	// Contains reports whether the bypass includes addr.
	IsWhitelist() bool
	Contains(ctx context.Context, network, addr string, opts ...Option) bool
}
