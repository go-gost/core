package ingress

import "context"

type GetOptions struct{}

type GetOption func(opts *GetOptions)

type SetOptions struct{}

type SetOption func(opts *SetOptions)

type Ingress interface {
	Get(ctx context.Context, host string, opts ...GetOption) string
	Set(ctx context.Context, host, endpoint string, opts ...SetOption)
}
