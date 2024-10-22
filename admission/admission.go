package admission

import "context"

type Options struct{}
type Option func(opts *Options)

type Admission interface {
	Admit(ctx context.Context, addr string, opts ...Option) bool
}
