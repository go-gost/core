package chain

import "context"

type SelectOptions struct {
	Addr string
}

type SelectOption func(*SelectOptions)

func AddrSelectOption(addr string) SelectOption {
	return func(o *SelectOptions) {
		o.Addr = addr
	}
}

type Hop interface {
	Nodes() []*Node
	Select(ctx context.Context, opts ...SelectOption) *Node
}
